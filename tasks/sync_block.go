package tasks

import (
	"context"
	"crynux_relay/blockchain"
	"crynux_relay/config"
	"crynux_relay/models"
	"errors"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func StartSyncBlockWithTerminateChannel(ch <-chan int) {

	syncedBlock, err := getSyncedBlock()

	if err != nil {
		log.Errorln("error getting synced block from the database")
		log.Fatal(err)
	}

	for {
		select {
		case stop := <-ch:
			if stop == 1 {
				return
			} else {
				processChannel(syncedBlock)
			}
		default:
			processChannel(syncedBlock)
		}
	}
}

func StartSyncBlock() {

	syncedBlock, err := getSyncedBlock()

	if err != nil {
		log.Errorln("error getting synced block from the database")
		log.Fatal(err)
	}

	for {
		processChannel(syncedBlock)
	}
}

func getSyncedBlock() (*models.SyncedBlock, error) {
	appConfig := config.GetConfig()
	syncedBlock := &models.SyncedBlock{}

	if err := config.GetDB().First(&syncedBlock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			syncedBlock.BlockNumber = appConfig.Blockchain.StartBlockNum
		} else {
			return nil, err
		}
	}

	return syncedBlock, nil
}

func processChannel(syncedBlock *models.SyncedBlock) {

	interval := 1
	batchSize := uint64(500)

	client, err := blockchain.GetRpcClient()
	if err != nil {
		log.Errorln("error getting the eth rpc client")
		log.Errorln(err)
		time.Sleep(time.Duration(interval) * time.Second)
		return
	}

	latestBlockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Errorln("error getting the latest block number")
		log.Errorln(err)
		time.Sleep(time.Duration(interval) * time.Second)
		return
	}

	if latestBlockNum <= syncedBlock.BlockNumber {
		time.Sleep(time.Duration(interval) * time.Second)
		return
	}

	log.Debugln("new block received: " + strconv.FormatUint(latestBlockNum, 10))

	for start := syncedBlock.BlockNumber + 1; start <= latestBlockNum; start += batchSize {

		end := start + batchSize - 1

		if end > latestBlockNum {
			end = latestBlockNum
		}

		log.Debugln("processing blocks from " +
			strconv.FormatUint(start, 10) +
			" to " +
			strconv.FormatUint(end, 10) +
			" / " +
			strconv.FormatUint(latestBlockNum, 10))

		if err := processTaskPending(start, end); err != nil {
			log.Errorf("processing task pending error: %v", err)
			time.Sleep(time.Duration(interval) * time.Second)
			return
		}

		if err := processTaskStarted(start, end); err != nil {
			log.Errorf("processing task started error: %v", err)
			time.Sleep(time.Duration(interval) * time.Second)
			return
		}

		if err := processTaskSuccess(start, end); err != nil {
			log.Errorf("processing task success error: %v", err)
			time.Sleep(time.Duration(interval) * time.Second)
			return
		}

		if err := processTaskAborted(start, end); err != nil {
			log.Errorf("processing task aborted error: %v", err)
			time.Sleep(time.Duration(interval) * time.Second)
			return
		}

		oldNum := syncedBlock.BlockNumber
		syncedBlock.BlockNumber = end
		if err := config.GetDB().Save(syncedBlock).Error; err != nil {
			syncedBlock.BlockNumber = oldNum
			log.Errorln(err)
			time.Sleep(time.Duration(interval) * time.Second)
		}

		if end != latestBlockNum {
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	time.Sleep(time.Duration(interval) * time.Second)
}

func processTaskPending(startBlockNum, endBlockNum uint64) error {
	taskContractInstance, err := blockchain.GetTaskContractInstance()
	if err != nil {
		return err
	}

	taskPendingEventIterator, err := taskContractInstance.FilterTaskPending(
		&bind.FilterOpts{
			Start:   startBlockNum,
			End:     &endBlockNum,
			Context: context.Background(),
		},
		nil,
	)

	if err != nil {
		return err
	}

	defer taskPendingEventIterator.Close()

	for {
		if !taskPendingEventIterator.Next() {
			break
		}

		taskPending := taskPendingEventIterator.Event

		log.Debugln("Task pending on chain: " +
			taskPending.TaskId.String() +
			"|" + taskPending.Creator.Hex() +
			"|" + string(taskPending.TaskHash[:]) +
			"|" + string(taskPending.DataHash[:]))

		task := &models.InferenceTask{}

		query := &models.InferenceTask{
			TaskId: taskPending.TaskId.Uint64(),
		}

		taskOnChain, err := blockchain.GetTaskById(taskPending.TaskId.Uint64())
		if err != nil {
			return err
		}

		attributes := &models.InferenceTask{
			Creator:   taskPending.Creator.Hex(),
			TaskHash:  hexutil.Encode(taskPending.TaskHash[:]),
			DataHash:  hexutil.Encode(taskPending.DataHash[:]),
			Status:    models.InferenceTaskCreatedOnChain,
			TaskType:  models.ChainTaskType(taskPending.TaskType.Int64()),
			VramLimit: taskOnChain.VramLimit.Uint64(),
		}

		if err := config.GetDB().Where(query).Attrs(attributes).FirstOrCreate(task).Error; err != nil {
			return err
		}
	}

	return nil
}

func processTaskStarted(startBlockNum, endBlockNum uint64) error {

	taskContractInstance, err := blockchain.GetTaskContractInstance()
	if err != nil {
		return err
	}

	taskStartedEventIterator, err := taskContractInstance.FilterTaskStarted(
		&bind.FilterOpts{
			Start:   startBlockNum,
			End:     &endBlockNum,
			Context: context.Background(),
		},
		nil,
		nil,
	)

	if err != nil {
		return err
	}

	defer taskStartedEventIterator.Close()

	taskStartedEvents := make(map[uint64][]models.SelectedNode)

	for {
		if !taskStartedEventIterator.Next() {
			break
		}

		taskStarted := taskStartedEventIterator.Event

		log.Debugln("Task created on chain: " +
			taskStarted.TaskId.String() +
			"|" + taskStarted.Creator.Hex() +
			"|" + string(taskStarted.TaskHash[:]) +
			"|" + string(taskStarted.DataHash[:]))

		taskId := taskStarted.TaskId.Uint64()
		taskStartedEvents[taskId] = append(taskStartedEvents[taskId], models.SelectedNode{NodeAddress: taskStarted.SelectedNode.Hex()})
	}

	for taskId, selectedNodes := range taskStartedEvents {
		task := &models.InferenceTask{TaskId: taskId}
		
		if err := config.GetDB().Where(task).First(task).Error; err != nil {
			return err
		}

		var existSelectedNodes []models.SelectedNode
		if err := config.GetDB().Model(task).Association("SelectedNodes").Find(&existSelectedNodes); err != nil {
			return err
		}
		if len(existSelectedNodes) == 0 {
			if err := config.GetDB().Model(task).Association("SelectedNodes").Append(selectedNodes); err != nil {
				return err
			}
		} else {
			existNodeAddresses := make(map[string]interface{})
			for _, node := range existSelectedNodes {
				existNodeAddresses[node.NodeAddress] = nil
			}
			
			var newSelectedNodes []models.SelectedNode
			for _, node := range selectedNodes {
				_, ok := existNodeAddresses[node.NodeAddress]
				if !ok {
					newSelectedNodes = append(newSelectedNodes, node)
				}
			}

			if err := config.GetDB().Model(task).Association("SelectedNodes").Append(newSelectedNodes); err != nil {
				return err
			}
		}

	}

	return nil
}

func processTaskSuccess(startBlockNum, endBlockNum uint64) error {
	taskContractInstance, err := blockchain.GetTaskContractInstance()
	if err != nil {
		return err
	}

	taskSuccessEventIterator, err := taskContractInstance.FilterTaskSuccess(
		&bind.FilterOpts{
			Start:   startBlockNum,
			End:     &endBlockNum,
			Context: context.Background(),
		},
		nil,
		nil,
	)

	if err != nil {
		return err
	}

	defer taskSuccessEventIterator.Close()

	for {
		if !taskSuccessEventIterator.Next() {
			break
		}

		taskSuccess := taskSuccessEventIterator.Event

		log.Debugln("Task success on chain: " +
			taskSuccess.TaskId.String() +
			"|" + string(taskSuccess.Result) +
			"|" + taskSuccess.ResultNode.Hex())

		task := &models.InferenceTask{
			TaskId: taskSuccess.TaskId.Uint64(),
		}

		if err := config.GetDB().Where(task).First(task).Error; err != nil {
			return err
		}

		if task.Status != models.InferenceTaskParamsUploaded {
			continue
		}

		selectedNode := &models.SelectedNode{
			InferenceTaskID: task.ID,
			NodeAddress:     taskSuccess.ResultNode.Hex(),
		}

		if err := config.GetDB().Where(selectedNode).First(selectedNode).Error; err != nil {
			return err
		}

		selectedNode.Result = hexutil.Encode(taskSuccess.Result)
		selectedNode.IsResultSelected = true

		if err := config.GetDB().Model(selectedNode).Select("Result", "IsResultSelected").Updates(selectedNode).Error; err != nil {
			return err
		}

		task.Status = models.InferenceTaskPendingResults

		if err := config.GetDB().Save(task).Error; err != nil {
			return err
		}
	}

	return nil
}

func processTaskAborted(startBlockNum, endBlockNum uint64) error {
	taskContractInstance, err := blockchain.GetTaskContractInstance()
	if err != nil {
		return err
	}

	taskAbortedEventIterator, err := taskContractInstance.FilterTaskAborted(
		&bind.FilterOpts{
			Start:   startBlockNum,
			End:     &endBlockNum,
			Context: context.Background(),
		},
		nil,
	)

	if err != nil {
		return err
	}

	defer taskAbortedEventIterator.Close()

	for {
		if !taskAbortedEventIterator.Next() {
			break
		}

		taskAborted := taskAbortedEventIterator.Event

		log.Debugln("Task aborted on chain: " + taskAborted.TaskId.String())

		task := &models.InferenceTask{
			TaskId: taskAborted.TaskId.Uint64(),
		}

		if err := config.GetDB().Where(task).First(task).Error; err != nil {
			return err
		}

		if task.Status == models.InferenceTaskResultsUploaded {
			continue
		}

		task.Status = models.InferenceTaskAborted

		if err := config.GetDB().Save(task).Error; err != nil {
			return err
		}
	}

	return nil
}
