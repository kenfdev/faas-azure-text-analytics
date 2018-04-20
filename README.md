# OpenFaaS Azure Text Analytics Function

This function is a simple wrapper function for the [Text Analytics API](https://azure.microsoft.com/en-us/services/cognitive-services/text-analytics/).

## Prerequisites

* [faas-cli](https://github.com/openfaas/faas-cli)
* [OpenFaaS](https://github.com/openfaas/faas) inside your cluster
* [API Key](https://azure.microsoft.com/en-us/try/cognitive-services/?api=text-analytics) for the Text Analytics API
* API Key deployed as a secret inside the cluster

## Try it out

### stack.yml

``` yaml
functions:
  azure-text-analytics:
    lang: golang-http
    handler: ./azure-text-analytics
    image: <Your DockerHub Account>/azure-text-analytics
    secrets:
    - azure-api-keys # Make sure you deployed the API Key as a secret with this name
    environment:
      write_debug: true
      API_KEY_NAME: azure-free-text-analytics-api-key # This is the Key of the secret
      TEXT_ANALYZE_BASE_URL: https://eastasia.api.cognitive.microsoft.com/text/analytics/v2.0
```

### Build, Push, Deploy

Be sure to set your OpenFaas Gateway URL with

* env variable `OPENFAAS_URL`
* `-g` option with the `faas-cli`

```
faas build && faas push && faas deploy [-g <GATEWAY URL>]
```

### Invoke

```
echo '{"text":"This function is an awesome function that simply integrates with the Azure Text Analytics API."}' | faas invoke azure-text-analytics [-g <GATEWAY URL>]
```

You'll get something like the following.

```json
{
  "language": "en",
  "sentiment_score": 0.884989857673645,
  "key_phrases": ["awesome function", "Azure Text Analytics API"]
}
```

## Give it a quick test with docker-compose

### API Key

Make sure you place your API key inside a file called `azure-free-text-analytics-api-key` .

Put your API key inside that file.

```
docker-compose up
```

By default, this docker-compose version will expose the API to http://127.0.0.1:9080.
Simply `curl` it with the `text` you want to analyze.

```
curl -v -X POST -d '{"text":"This function is an awesome function that simply integrates with the Azure Text Analytics API."}' http://localhost:9080
```
