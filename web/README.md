## Cloud Run

### Cloud Build CI/CD

#### Build Image
```shell
docker build -t gcr.io/clouddeveloper-299318/bitbucket.org/sachindramaharjan/service/web$SHORT_SHA .
```

#### Push Image to Container Registry
```shell
docker push gcr.io/clouddeveloper-299318/bitbucket.org/sachindramaharjan/service/web:$SHORT_SHA
```

#### Deploy to Cloud Run

``` shell
gcloud run deploy web-project --region=us-central1 --platform=managed --image=gcr.io/clouddeveloper-299318/bitbucket.org/sachindramaharjan/service/web:SHORT_SHA
```
