pipeline {
    agent {
        label 'docker'
    }

    options {
        ansiColor('xterm')
        disableConcurrentBuilds()
    }

    environment {
        DOCKER_REGISTRY = "docker-registry.dc1.apstra.com:5000"
        TAG = "${DOCKER_REGISTRY}/goapstra-sources:${GIT_COMMIT}"
    }

    stages {
        stage('Build docker image with source code') {
            steps {
                sh '''
                    docker build \
                        --label "jenkins/BRANCH_NAME=${BRANCH_NAME}" \
                        --label "jenkins/BUILD_NUMBER=${BUILD_NUMBER}" \
                        --label "jenkins/BUILD_URL=${BUILD_URL}" \
                        --label "jenkins/GIT_COMMIT=${GIT_COMMIT}" \
                        --label "jenkins/JOB_NAME=${JOB_NAME}" \
                        --label "jenkins/JOB_URL=${JOB_URL}" \
                        --label "repo=git@bitbucket.org:apstrktr/goapstra.git" \
                        --tag ${TAG} \
                        --file ci.Dockerfile \
                        . \
                    && docker push ${TAG}
                '''
            }
        }
        stage('Run tests') {
            parallel {
                stage('make lint-revive') {
                    steps {
                        sh 'docker run --rm ${TAG} make lint-revive'
                    }
                }
                stage('make lint-staticcheck') {
                    steps {
                        sh 'docker run --rm ${TAG} make lint-staticcheck'
                    }
                }
                stage('make fmt-check') {
                    steps {
                        sh 'docker run --rm ${TAG} make fmt-check'
                    }
                }
                stage('make vet') {
                    steps {
                        sh 'docker run --rm ${TAG} make vet'
                    }
                }
                stage('make unit-tests') {
                    steps {
                        sh 'docker run --rm ${TAG} make unit-tests'
                    }
                }
            }
        }
    }
    post {
        cleanup {
            sh 'docker image rm -f ${TAG}'
        }
    }
}
