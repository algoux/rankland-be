// model 存放数据模型定义的目录
// 将该应用所需的所有 db 表映射的模型存放在该目录下，
// 定义的 model 为了支持能够在迁移机器时正常使用，需要在
// struct tag 中写明每个字段对应的 db 约束
package model