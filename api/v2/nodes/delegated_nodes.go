package nodes

import (
	"crynux_relay/api/v1/response"
	"crynux_relay/config"
	"crynux_relay/models"

	"github.com/gin-gonic/gin"
)

type GetDelegatedNodesInput struct {
	Page     *int `json:"page" query:"page" description:"page number" validate:"required,min=1" default:"1"`
	PageSize *int `json:"page_size" query:"page_size" description:"number of items per page" validate:"required,min=1,max=100" default:"30"`
}

type GetDelegatedNodesOutput struct {
	response.Response
	Data []*Node `json:"data"`
}

func GetDelegatedNodes(c *gin.Context, input *GetDelegatedNodesInput) (*GetDelegatedNodesOutput, error) {
	page := 1
	if input.Page != nil {
		page = *input.Page
	}
	pageSize := 30
	if input.PageSize != nil {
		pageSize = *input.PageSize
	}

	nodes, err := models.GetDelegatedNodes(c.Request.Context(), config.GetDB(), (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, response.NewExceptionResponse(err)
	}

	nodeDatas := make([]*Node, 0)
	errCh := make(chan error, len(nodes))
	semaphone := make(chan struct{}, 10)
	for _, node := range nodes {
		go func(n *models.Node) {
			semaphone <- struct{}{}
			defer func() {
				<-semaphone
			}()
			nodeData, err := getNodeData(c.Request.Context(), n)
			if err != nil {
				errCh <- err
				return
			}
			nodeDatas = append(nodeDatas, nodeData)
		}(node)
	}
	for i := 0; i < len(nodes); i++ {
		if err := <-errCh; err != nil {
			return nil, response.NewExceptionResponse(err)
		}
	}

	return &GetDelegatedNodesOutput{
		Data: nodeDatas,
	}, nil
}
