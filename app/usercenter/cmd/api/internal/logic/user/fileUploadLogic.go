package user

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"
	"homestay/common/ctxdata"
	"net/http"
	"os"
	"path/filepath"
)

const maxFileSize = 10 << 20 // 10 MB

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//	func (l *FileUploadLogic) FileUpload(req *types.FileUploadReq) (resp *types.FileUploadResp, err error) {
//		// todo: add your logic here and delete this line
//		// 数据入库
//		rp := &model.RepositoryPool{
//			Identity: helper.UUID(),
//			Hash:     req.Hash,
//			Name:     req.Name,
//			Ext:      req.Ext,
//			Size:     req.Size,
//			Path:     req.Path,
//		}
//		_, err = l.svcCtx.Engine.Insert(rp)
//		if err != nil {
//			return
//		}
//
//		resp = new(types.FileUploadReply)
//		resp.Identity = rp.Identity
//		resp.Ext = rp.Ext
//		resp.Name = rp.Name
//
//		return
//	}
func (l *FileUploadLogic) FileUpload(r *http.Request) (resp *types.FileUploadResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	_ = userId
	_ = r.ParseMultipartForm(maxFileSize)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	//fmt.Println("GetAccessKeyID: " + l.svcCtx.Config.OSS_ACCESS_KEY_ID)
	//fmt.Println("GetAccessKeySecret: " + l.svcCtx.Config.OSS_ACCESS_KEY_SECRET)

	//fmt.Println("GetAccessKeyID: " + os.Getenv("OSS_ACCESS_KEY_ID"))
	//fmt.Println("GetAccessKeySecret: " + os.Getenv("OSS_ACCESS_KEY_SECRET"))
	fileName := uuid.New().String() + filepath.Ext(handler.Filename)
	//fmt.Println("fileName: " + fileName)

	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	//fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	//fmt.Println("GetAccessKeyID: " + provider.GetCredentials().GetAccessKeyID())
	//fmt.Println("GetAccessKeySecret: " + provider.GetCredentials().GetAccessKeySecret())
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com",
		"",
		"",
		oss.SetCredentialsProvider(&provider))

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	bucketName := "web-sh-ti"
	// 填写存储空间名称
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	err = bucket.PutObject(fileName, file)
	if err != nil {
		return nil, err
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	//err = bucket.PutObjectFromFile("exampledir/exampleobject.txt", "D:\\localpath\\examplefile.txt")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}

	//tempFile, err := os.Create(path.Join(l.svcCtx.Config.Path, handler.Filename))
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	//defer tempFile.Close()
	//io.Copy(tempFile, file)
	url := "https://" + bucketName + "." + "oss-cn-hangzhou.aliyuncs.com" + "/" + fileName

	return &types.FileUploadResp{
		URL: url,
	}, nil
}
