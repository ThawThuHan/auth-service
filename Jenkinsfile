pipeline {
    agent none

    environment {
        APP_NAME = "auth-service"
        HARBOR_REGISTRY = "10.11.0.20"
        HARBOR_PROJECT = "kubernetes-testing"
        HELM_CHART = "auth-service"
    }

    stages {
        stage('Prepare') {
            steps {
                script {
                    if (env.CHANGE_ID) {
                        echo "This is a pull request build for PR #${env.CHANGE_ID}"
                        env.APP_VERSION = "pr-${env.CHANGE_ID}-${env.BUILD_NUMBER}"
                        env.FULL_IMAGE_TAG = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}:${env.APP_VERSION}"
                        env.HARBOR_ROBOT_CREDENTIAL_ID = "harbor-robot-for-kubernetes-testing"
                    } else {
                        echo "This is a branch build for ${env.BRANCH_NAME}"
                        env.APP_VERSION = "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
                        env.FULL_IMAGE_TAG = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}:${env.APP_VERSION}"
                        env.LATEST_IMAGE_TAG = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}:latest"
                        env.HARBOR_ROBOT_CREDENTIAL_ID = "harbor-robot-for-kubernetes-prod"
                    }
                }
            }
        }

        stage('Checkout') {
            agent any
            steps {
                checkout scm
            }
        }

        stage('Run Tests') {
            agent {
                docker {
                    image 'golang:1.25.1'
                    args '-v $HOME/go/pkg/mod:/go/pkg/mod'
                }
            }

            steps {
                sh '''
                    go mod tidy
                    go test ./...
                '''
            }
        }

        stage('Build and Push Docker Image') {
            agent any
            steps {
                script {
                    withCredentials([
                        usernamePassword(
                            credentialsId: env.HARBOR_ROBOT_CREDENTIAL_ID, 
                            usernameVariable: 'HARBOR_USER', 
                            passwordVariable: 'HARBOR_PASSWORD'
                        )]) {
                        sh '''
                            docker login ${HARBOR_REGISTRY} -u ${HARBOR_USER} -p ${HARBOR_PASSWORD}
                            docker build -t ${FULL_IMAGE_TAG} .
                            docker push ${FULL_IMAGE_TAG}
                            if [ -n "${LATEST_IMAGE_TAG}" ]; then
                                docker tag ${FULL_IMAGE_TAG} ${LATEST_IMAGE_TAG}
                                docker push ${LATEST_IMAGE_TAG}
                            fi
                        '''
                    }
                }
            }
        }

        stage('Build Helm Chart') {
            agent {
                docker {
                    image 'alpine/helm:4'
                }
            }
            steps {
                sh '''
                    echo "helm chart building..."
                    helm package ./helm-chart --version ${APP_VERSION} --app-version ${APP_VERSION} -d ./helm-chart-output
                '''
                stash includes: 'helm-chart-output/*.tgz', name: 'helm-chart'
            }
        }

        stage('Push Helm Chart to Harbor') {
            agent {
                docker {
                    image 'alpine/helm:4'
                }
            }
            steps {
                script {
                    withCredentials([
                        usernamePassword(
                            credentialsId: env.HARBOR_ROBOT_CREDENTIAL_ID, 
                            usernameVariable: 'HARBOR_USER', 
                            passwordVariable: 'HARBOR_PASSWORD'
                        )]) {
                        unstash 'helm-chart'
                        sh '''
                            echo "helm chart pushing..."
                            helm registry login -u ${HARBOR_USER} --password-stdin ${HARBOR_REGISTRY} <<< ${HARBOR_PASSWORD}
                            helm push helm-chart-output/*.tgz oci://${HARBOR_REGISTRY}/${HARBOR_PROJECT}
                        '''
                    }
                }
            }
        }

        stage('Deploy') {
            agent {
                docker {
                    image 'alpine/helm:latest'
                    args '--entrypoint=""'
                }
            }
            steps {
                withCredentials([
                    file(credentialsId: 'Jenkins-KubeConfig', variable: 'KUBECONFIG'),
                    file(credentialsId: 'PRIVATE_KEY', variable: 'PRIVATE_KEY'),
                    file(credentialsId: 'PUBLIC_KEY', variable: 'PUBLIC_KEY'),
                    string(credentialsId: 'JWT_SECRET', variable: 'JWT_SECRET'),
                    usernamePassword(
                        credentialsId: env.HARBOR_ROBOT_CREDENTIAL_ID, 
                        usernameVariable: 'HARBOR_USER', 
                        passwordVariable: 'HARBOR_PASSWORD'
                    )
                ]) {
                sh '''
                    export KUBECONFIG=$KUBECONFIG
                    helm registry login -u ${HARBOR_USER} --password-stdin ${HARBOR_REGISTRY} <<< ${HARBOR_PASSWORD}
                    helm upgrade --install ${APP_NAME} oci://${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${HELM_CHART} \
                        --namespace testing --create-namespace \
                        --set image.repository=${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME} \
                        --set image.tag=${APP_VERSION} \
                        --set-file secret.data.PRIVATE_KEY=${PRIVATE_KEY} \
                        --set-file config.PUBLIC_KEY=${PUBLIC_KEY} \
                        --set secret.data.JWT_SECRET=${JWT_SECRET} \
                        --wait --timeout 5m0s
                '''
                }
            }
        }
    }
}