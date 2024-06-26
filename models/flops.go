package models

import "gorm.io/gorm"

type NetworkFLOPS struct {
	gorm.Model
	GFLOPS float64 `json:"gflops"`
}

var gpuGFLOPSMap map[string]float64 = map[string]float64{
	"NVIDIA GeForce GTX 1050 Ti": 2138.1,
	"NVIDIA GeForce GTX 1060": 4375.0,
	"NVIDIA GeForce GTX 1070": 6462.7,
	"NVIDIA GeForce GTX 1070 Ti": 8186.1,
	"NVIDIA GeForce GTX 1080": 8872.9,
	"NVIDIA GeForce GTX 1080 Ti": 11339.7,
	"NVIDIA TITAN X Pascal": 10974.2,
	"NVIDIA TITAN Xp": 12149.7,
	"NVIDIA TITAN V": 14899.2,
	"NVIDIA GeForce GTX 1650": 2849.28,
	"NVIDIA GeForce GTX 1650 Super": 4416.00,
	"NVIDIA GeForce GTX 1660": 5027.00,
	"NVIDIA GeForce GTX 1660 Super": 5027.00,
	"NVIDIA GeForce GTX 1660 Ti": 5437.44,
	"NVIDIA GeForce RTX 2060": 6451.20,
	"NVIDIA GeForce RTX 2060 Super": 7180.00,
	"NVIDIA GeForce RTX 2070": 7464.96,
	"NVIDIA GeForce RTX 2070 Super": 9060.00,
	"NVIDIA GeForce RTX 2080": 10068.48,
	"NVIDIA GeForce RTX 2080 Super": 11150.00,
	"NVIDIA GeForce RTX 2080 Ti": 13447.68,
	"NVIDIA TITAN RTX": 16312.32,
	"NVIDIA GeForce RTX 3050": 9.01 * 1024,
	"NVIDIA GeForce RTX 3060": 12.74 * 1024,
	"NVIDIA GeForce RTX 3060 Ti": 16.20 * 1024,
	"NVIDIA GeForce RTX 3070": 20.31 * 1024,
	"NVIDIA GeForce RTX 3070 Ti": 21.75 * 1024,
	"NVIDIA GeForce RTX 3080": 30.6 * 1024,
	"NVIDIA GeForce RTX 3080 Ti": 34.1 * 1024,
	"NVIDIA GeForce RTX 3090": 35.58 * 1024,
	"NVIDIA GeForce RTX 3090 Ti": 40 * 1024,
	"NVIDIA GeForce RTX 4060": 15.1 * 1024,
	"NVIDIA GeForce RTX 4060 Ti": 22.1 * 1024,
	"NVIDIA GeForce RTX 4070": 29.1 * 1024,
	"NVIDIA GeForce RTX 4070 Super": 35.48 * 1024,
	"NVIDIA GeForce RTX 4070 Ti": 40.1 * 1024,
	"NVIDIA GeForce RTX 4070 Ti Super": 44.10 * 1024,
	"NVIDIA GeForce RTX 4080": 48.7 * 1024,
	"NVIDIA GeForce RTX 4080 Super": 52.22 * 1024,
	"NVIDIA GeForce RTX 4090 D": 73.5 * 1024,
	"NVIDIA GeForce RTX 4090": 82.6 * 1024,
	"NVIDIA RTX A4000": 19170,
	"NVIDIA RTX A4500": 23655.91,
	"NVIDIA RTX A5000": 27772.65,
	"NVIDIA RTX A5500": 34101.38,
	"NVIDIA RTX A6000": 38709.67,
	"NVIDIA RTX 4000 Ada": 26730,
	"NVIDIA RTX 4500 Ada": 39630,
	"NVIDIA RTX 5000 Ada": 65280,
	"NVIDIA RTX 5880 Ada": 69272,
	"NVIDIA RTX 6000 Ada": 91060,
	"NVIDIA A2": 4531,
	"NVIDIA A10": 31240,
	"NVIDIA A16": 4608,
	"NVIDIA A30": 10320,
	"NVIDIA A40": 37420,
	"NVIDIA A100": 19500,
	"NVIDIA H100": 51200,
	"NVIDIA L40": 90516,
	"NVIDIA L4": 30300,
	"NVIDIA GeForce RTX 2060 Laptop GPU": 4608,
	"NVIDIA GeForce RTX 2070 Laptop GPU": 6636,
	"NVIDIA GeForce RTX 2070 Super Laptop GPU": 7066,
	"NVIDIA GeForce RTX 2080 Laptop GPU": 9362,
	"NVIDIA GeForce RTX 2080 Super Laptop GPU": 9585,
	"NVIDIA GeForce RTX 3050 Laptop GPU": 7.13 * 1024,
	"NVIDIA GeForce RTX 3060 Laptop GPU": 13.07 * 1024,
	"NVIDIA GeForce RTX 3070 Laptop GPU": 16.59 * 1024,
	"NVIDIA GeForce RTX 3070 Ti Laptop GPU": 17.49 * 1024,
	"NVIDIA GeForce RTX 3080 Laptop GPU": 21.01 * 1024,
	"NVIDIA GeForce RTX 3080 Ti Laptop GPU": 23.60 * 1024,
	"NVIDIA GeForce RTX 4050 Laptop GPU": 12.1 * 1024,
	"NVIDIA GeForce RTX 4060 Laptop GPU": 14.5 * 1024,
	"NVIDIA GeForce RTX 4070 Laptop GPU": 20.0 * 1024,
	"NVIDIA GeForce RTX 4080 Laptop GPU": 33.8 * 1024,
	"NVIDIA GeForce RTX 4090 Laptop GPU": 50.1 * 1024,
	"NVIDIA RTX 2000 Ada Generation Laptop GPU": 14500,
	"NVIDIA RTX 3000 Ada Generation Laptop GPU": 19900,
	"NVIDIA RTX 3500 Ada Generation Laptop GPU": 23300,
	"NVIDIA RTX 4000 Ada Generation Laptop GPU": 33600,
	"NVIDIA RTX 5000 Ada Generation Laptop GPU": 42600,
	"Apple M1 Type": 2.6 * 1024,
	"Apple M1 Max Type": 10.4 * 1024,
	"Apple M1 Pro Type": 5.3 * 1024,
	"Apple M1 Ultra Type": 21 * 1024,
	"Apple M2 Type": 3.6 * 1024,
	"Apple M2 Max Type": 13.6 * 1024,
	"Apple M2 Pro Type": 6.8 * 1024,
	"Apple M2 Ultra Type": 27.2 * 1024,
	"Apple M3 Type": 4.1 * 1024,
	"Apple M3 Max Type": 14.131 * 1024,
	"Apple M3 Pro Type": 6.359 * 1024,
}


func GetGPUGFLOPS(gpuName string) float64 {
	if gflops, ok := gpuGFLOPSMap[gpuName]; ok {
		return gflops
	} else {
		return 10 * 1024
	}
}