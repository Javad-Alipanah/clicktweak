ClickTweak [^1]
----------

# ClickTweak is a fast URL shortener

# Installation Guide

#### Prequsites: `docker docker-compose`

1. clone repository:
    ```bash
    git clone https://gitlab.com/javadalipanah/clicktweak.git && cd clicktweak
    ```
2. to prepare service images, you have two alternatives:
   * login to my private docker registry to be able to pull images
        ```bash
        dokcer login -u ${REPO_USER} -p ${REPO_PASSWORD}
        ```
   * build images:
        ```bash
        make static && \
        docker build -t reg.alipanah.me/core -f build/package/core.Dockerfile . \
        docker build -t reg.alipanah.me/dispatcher -f build/package/dispatcher.Dockerfile . \
        docker build -t reg.alipanah.me/consumer -f build/package/consumer.Dockerfile . \
        docker build -t reg.alipanah.me/analyzer -f build/package/analyzer.Dockerfile .
        ```
3. run docker-compose:
    ```bash
    docker-compose -f deployments/docker-compose.yml down
    ```

4. check service health:
    * `Core: localhost:8080`
    * `Dispatcher: localhost:8081`
    * `Analyzer: localhost:8082`
___

[^1]: This project is an implementation for the [Yektanet](https://en.yektanet.com/)'s interview project. 