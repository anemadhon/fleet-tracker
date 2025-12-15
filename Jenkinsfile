pipeline {
    agent any
    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '20'))
    }
    environment {
        COMPOSE_PROJECT_NAME = "fleet-${env.BUILD_NUMBER}"
        DOCKER_BUILDKIT = '1'
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Prepare Environment') {
            steps {
                script {
                    // Copy .env.example to .env if exists
                    bat '''
                        if exist .env.example (
                            copy .env.example .env
                        ) else (
                            if not exist .env (
                                type nul > .env
                            )
                        )
                    '''
                }
            }
        }
        stage('Build Docker Images') {
            steps {
                bat 'docker-compose build --no-cache'
            }
        }
        stage('Start Services') {
            steps {
                bat 'docker-compose up -d'
                // bat 'docker-compose up -d postgres rabbitmq mosquitto'
                // bat 'timeout /t 10 /nobreak' // Wait 10 seconds
            }
        }
        stage('Run Tests') {
            steps {
                script {
                    // Run tests inside a temporary Go container
                    bat '''
                        docker run --rm ^
                            --network %COMPOSE_PROJECT_NAME%_fleet-network ^
                            -v %cd%:/app ^
                            -w /app ^
                            -e CGO_ENABLED=0 ^
                            -e GO111MODULE=on ^
                            golang:1.21 ^
                            sh -c "go mod tidy && go test ./... -v"
                    '''
                }
            }
        }
        stage('Build & Start All Services') {
            steps {
                bat 'docker-compose up -d'
                // bat 'timeout /t 5 /nobreak'
            }
        }
        stage('Health Check') {
            steps {
                script {
                    bat '''
                        echo Checking services status...
                        docker-compose ps
                    '''
                    
                    // Optional: API health check with curl (needs curl installed)
                    // bat 'curl -f http://localhost:8093/health || echo API health check skipped'
                }
            }
        }
        stage('Package Docker Images') {
            when {
                branch 'main'
            }
            steps {
                script {
                    bat '''
                        docker save -o api-image.tar %COMPOSE_PROJECT_NAME%_api:latest
                        docker save -o subscriber-image.tar %COMPOSE_PROJECT_NAME%_subscriber:latest
                        docker save -o worker-image.tar %COMPOSE_PROJECT_NAME%_worker:latest
                        docker save -o publisher-image.tar %COMPOSE_PROJECT_NAME%_publisher:latest
                        
                        tar -czf docker-images.tar.gz api-image.tar subscriber-image.tar worker-image.tar publisher-image.tar
                        del *-image.tar
                    '''
                    
                    archiveArtifacts artifacts: 'docker-images.tar.gz', fingerprint: true
                }
            }
        }
        stage('Deploy (placeholder)') {
            when {
                branch 'main'
            }
            steps {
                echo 'Deploy your containerized services here'
            }
        }
    }
    post {
        always {
            script {
                // Stop and remove containers
                bat 'docker-compose down -v || exit 0'
                
                // Clean up dangling images
                bat 'docker image prune -f || exit 0'
            }
        }
        success {
            echo "Pipeline succeeded for commit ${env.GIT_COMMIT}"
        }
        failure {
            echo "Pipeline failed for commit ${env.GIT_COMMIT}"
            script {
                // Show logs for debugging
                bat 'docker-compose logs --tail=100 || exit 0'
            }
        }
    }
}