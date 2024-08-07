package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
)

type FileHead struct {
	Key          string
	LastModified time.Time
	NumBytes     int64
}

type DownloadFileByKeyOutput struct {
	Bytes    []byte
	NumBytes int64
}

type Bucket interface {
	GetAllFileHeads(ctx context.Context, prefix string) ([]*FileHead, error)
	DownloadFileByKey(ctx context.Context, key string) (*DownloadFileByKeyOutput, error)
	UploadLargeFile(ctx context.Context, key string, file io.Reader) error
	UploadFile(ctx context.Context, key string, file io.Reader) error
}

type bucket struct {
	name                         string
	fileDownloadTL, fileUploadTL time.Duration

	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
}

func NewBucket(bucketName string, fileDownloadTL time.Duration, fileUploadTL time.Duration) (Bucket, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("error bucketName cant be empty")
	}
	log.Infoln(fileDownloadTL, fileUploadTL)
	//if fileDownloadTL.Seconds() < 1 || fileDownloadTL.Seconds() > 1000 {
	//  return nil, fmt.Errorf("error fileDownloadTL must be in range 1 - 40 seconds")
	//}
	//if fileUploadTL.Seconds() < 1 || fileUploadTL.Seconds() > 40 {
	//  return nil, fmt.Errorf("error fileUploadTL must be in range 1 - 40 seconds")
	//}
	// Создаем кастомный обработчик эндпоинтов, который для сервиса S3 и региона ru-central1 выдаст корректный URL
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: "ru-central1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Подгружаем конфигрурацию из ~/.aws/*
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}
	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	// Запрашиваем список бакетов
	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range result.Buckets {
		log.Printf("bucket=%s creation time=%s", aws.ToString(bucket.Name), bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}

	var partMBytes int64 = 4
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = partMBytes * 1024 * 1024
	})

	downloader := manager.NewDownloader(client)

	return &bucket{
		name:           bucketName,
		fileDownloadTL: fileDownloadTL,
		fileUploadTL:   fileUploadTL,

		client:     client,
		uploader:   uploader,
		downloader: downloader,
	}, err
}

func NewUniversalBucket(
	s3URL, region, bucketName, partitionID string,
	fileDownloadTL time.Duration, fileUploadTL time.Duration,
) (Bucket, error) {
	if s3URL == "" {
		return nil, fmt.Errorf("error s3 url cant be empty")
	} else if _, err := url.Parse(s3URL); err != nil {
		return nil, fmt.Errorf("error parse s3 url: %w", err)
	}
	if region == "" {
		return nil, fmt.Errorf("error s3 region cant be empty")
	}
	if bucketName == "" {
		return nil, fmt.Errorf("error s3 bucketName cant be empty")
	}
	if partitionID == "" {
		return nil, fmt.Errorf("error s3 partitionID cant be empty")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:   partitionID,
				URL:           s3URL,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// load config from ~/.aws/*
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	//crate s3 client
	client := s3.NewFromConfig(cfg)

	// ping bucket
	_, err = client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("error while ListBuckets: %w", err)
	}

	var partMBytes int64 = 4
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = partMBytes * 1024 * 1024
	})

	downloader := manager.NewDownloader(client)

	return &bucket{
		name:           bucketName,
		fileDownloadTL: fileDownloadTL,
		fileUploadTL:   fileUploadTL,

		client:     client,
		uploader:   uploader,
		downloader: downloader,
	}, err
}

func NewBucketWithCredentials(
	ctx context.Context,
	keyID,
	accessKey,
	region,
	s3url,
	bucketName string,
	fileDownloadTL time.Duration, fileUploadTL time.Duration,
) (Bucket, error) {
	if keyID == "" {
		return nil, fmt.Errorf("error s3 keyID cant be empty")
	}
	if accessKey == "" {
		return nil, fmt.Errorf("error s3 accessKey cant be empty")
	}
	if region == "" {
		return nil, fmt.Errorf("error s3 region cant be empty")
	}
	if s3url == "" {
		return nil, fmt.Errorf("error s3 url cant be empty")
	} else if _, err := url.Parse(s3url); err != nil {
		return nil, fmt.Errorf("error parse s3 url: %w", err)
	}
	if bucketName == "" {
		return nil, fmt.Errorf("error s3 bucketName cant be empty")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:           s3url,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(keyID, accessKey, ""),
		),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	//ping
	_, err = client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("error while ListBuckets: %w", err)
	}

	return &bucket{
		name:           bucketName,
		fileDownloadTL: fileDownloadTL,
		fileUploadTL:   fileUploadTL,

		client: client,
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 4 * 1024 * 1024
		}),
		downloader: manager.NewDownloader(client),
	}, err
}

func (b *bucket) UploadFile(ctx context.Context, key string, file io.Reader) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (b *bucket) UploadLargeFile(ctx context.Context, key string, file io.Reader) error {
	err := bucketExists(ctx, b.client, b.name)
	if err != nil {
		return err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*b.fileUploadTL)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	_, err = b.uploader.Upload(ctxTimeout, &s3.PutObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (b *bucket) listFiles(ctx context.Context, prefix string) ([]types.Object, error) {
	paginator := s3.NewListObjectsV2Paginator(
		b.client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(b.name),
			Prefix: aws.String(prefix),
		},
		func(o *s3.ListObjectsV2PaginatorOptions) {
			o.Limit = 512
		},
	)
	contents := make([]types.Object, 0)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		contents = append(contents, page.Contents...)
	}

	return contents, nil
}

func (b *bucket) GetAllFileHeads(ctx context.Context, prefix string) ([]*FileHead, error) {
	res, err := b.listFiles(ctx, prefix)
	if err != nil {
		return nil, err
	}

	var fhs []*FileHead
	for _, c := range res {
		fhs = append(fhs, &FileHead{Key: *c.Key, LastModified: *c.LastModified, NumBytes: c.Size})
	}

	return fhs, nil
}

// DownloadFileByKey locks goroutine until the end of download or ctxDw or ctx
func (b *bucket) DownloadFileByKey(ctx context.Context, key string) (*DownloadFileByKeyOutput, error) {
	obj, err := b.GetFileMetadata(ctx, key)
	if err != nil {
		return &DownloadFileByKeyOutput{nil, 0}, err
	}

	buf := make([]byte, obj.ContentLength)
	wr := manager.NewWriteAtBuffer(buf)
	log.Debugf("Downloading file %s\n", key)
	ctxDw, cancel := context.WithTimeout(ctx, b.fileDownloadTL)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	numBytes, err := b.downloader.Download(
		ctxDw,
		wr,
		&s3.GetObjectInput{
			Bucket: aws.String(b.name),
			Key:    aws.String(key),
		},
		func(downloader *manager.Downloader) {
			downloader.Concurrency = 3
		},
	)

	if err != nil {
		return &DownloadFileByKeyOutput{nil, 0}, err
	}

	return &DownloadFileByKeyOutput{buf, numBytes}, nil
}

func (b *bucket) GetFileMetadata(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	object, err := b.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func bucketExists(ctx context.Context, client *s3.Client, name string) error {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	return err
}
