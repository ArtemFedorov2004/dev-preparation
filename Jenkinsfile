pipeline {
    agent any

    environment {
        REPO_URL        = 'https://github.com/ArtemFedorov2004/dev-preparation'
        BRANCH          = 'main'

        VPS_USER        = 'root'
        VPS_DEPLOY_DIR  = '/opt/devprep'

        REGISTRY        = ''

        FRONTEND_VERSION = '0.1.0-SNAPSHOT'
        FRONTEND_JAR     = "target/devprep-frontend-${FRONTEND_VERSION}.jar"

        IMAGE_BACKEND    = 'devprep/backend'
        IMAGE_FRONTEND   = 'devprep/frontend'
        IMAGE_TAG        = "${BUILD_NUMBER}"
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 40, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: env.BRANCH,
                    url: env.REPO_URL
            }
        }

        stage('Build Frontend') {
            steps {
                withMaven(maven: 'Maven 3') {
                    dir('frontend') {
                        sh 'mvn clean package'
                    }
                }
            }
        }

        stage('Build Backend') {
            steps {
                dir('backend') {
                    sh 'go build ./cmd/api/main.go'
                }
            }
        }

        stage('Build Docker Images') {
            steps {
                withCredentials([string(credentialsId: 'vps-host', variable: 'VPS_HOST')]) {
                    sh """
                        docker build \
                            -t \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:${env.IMAGE_TAG} \
                            -t \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:latest \
                            ./backend

                        docker build \
                            --build-arg JAR_FILE=${env.FRONTEND_JAR} \
                            -t \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:${env.IMAGE_TAG} \
                            -t \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:latest \
                            ./frontend
                    """
                }
            }
        }

        stage('Ensure Registry') {
            steps {
                withCredentials([
                    sshUserPrivateKey(credentialsId: 'vps-ssh-key', keyFileVariable: 'SSH_KEY'),
                    string(credentialsId: 'vps-host', variable: 'VPS_HOST')
                ]) {
                    sh """
                        ssh -i \$SSH_KEY -o StrictHostKeyChecking=no ${env.VPS_USER}@\${VPS_HOST} '
                            docker ps --format "{{.Names}}" | grep -q "^registry\$" || \
                            docker run -d \
                                --name registry \
                                --restart always \
                                -p 5000:5000 \
                                -v /opt/registry:/var/lib/registry \
                                registry:2
                        '
                    """
                }
            }
        }

        stage('Push Images') {
            steps {
                withCredentials([string(credentialsId: 'vps-host', variable: 'VPS_HOST')]) {
                    sh """
                        docker push \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:${env.IMAGE_TAG}
                        docker push \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:latest
                        docker push \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:${env.IMAGE_TAG}
                        docker push \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:latest
                    """
                }
            }
        }

        stage('Deploy') {
            steps {
                withCredentials([
                    sshUserPrivateKey(
                        credentialsId: 'vps-ssh-key',
                        keyFileVariable: 'SSH_KEY'
                    ),
                    string(credentialsId: 'vps-host', variable: 'VPS_HOST'),
                    file(
                        credentialsId: 'devprep-env-file',
                        variable: 'ENV_FILE'
                    )
                ]) {
                    sh """
                        SSH="ssh -i \$SSH_KEY -o StrictHostKeyChecking=no ${env.VPS_USER}@\${VPS_HOST}"
                        SCP="scp -i \$SSH_KEY -o StrictHostKeyChecking=no"

                        \$SSH 'mkdir -p ${env.VPS_DEPLOY_DIR}'

                        \$SCP infra/docker-compose.yml infra/docker-compose.prod.yml \
                            ${env.VPS_USER}@\${VPS_HOST}:${env.VPS_DEPLOY_DIR}

                        \$SSH 'mkdir -p ${env.VPS_DEPLOY_DIR}/{postgres,keycloak,nginx}'

                        \$SCP -r infra/postgres infra/nginx \
                            ${env.VPS_USER}@\${VPS_HOST}:${env.VPS_DEPLOY_DIR}

                        \$SCP infra/keycloak/realm-export.prod.json \
                            ${env.VPS_USER}@\${VPS_HOST}:${env.VPS_DEPLOY_DIR}/keycloak/realm-export.json

                        \$SCP \$ENV_FILE \
                            ${env.VPS_USER}@\${VPS_HOST}:${env.VPS_DEPLOY_DIR}/.env

                        \$SSH '
                            cd ${env.VPS_DEPLOY_DIR}
                            docker compose -f docker-compose.yml -f docker-compose.prod.yml pull backend frontend
                            docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --remove-orphans
                            docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
                        '
                    """
                }
            }
        }
    }

    post {
        success {
            echo "✅ Build #${BUILD_NUMBER} успешно задеплоен!"
        }
        failure {
            echo "❌ Build #${BUILD_NUMBER} провалился. Проверьте логи."
        }
        always {
            withCredentials([string(credentialsId: 'vps-host', variable: 'VPS_HOST')]) {
                sh """
                    docker rmi \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:${env.IMAGE_TAG} || true
                    docker rmi \${VPS_HOST}:5000/${env.IMAGE_BACKEND}:latest || true
                    docker rmi \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:${env.IMAGE_TAG} || true
                    docker rmi \${VPS_HOST}:5000/${env.IMAGE_FRONTEND}:latest || true
                """
            }
            cleanWs()
        }
    }
}