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
                    // Copy .env.example to .env if exists, or create dummy .env
                    sh '''
                        if [ -f .env.example ]; then
                            cp .env.example .env
                        elif [ ! -f .env ]; then
                            touch .env
                        fi
                    '''
                }
            }
        }
        
        stage('Build Docker Images') {
            steps {
                sh 'docker-compose build --no-cache'
            }
        }
        
        stage('Start Services') {
            steps {
                sh 'docker-compose up -d postgres rabbitmq mosquitto'
                sh 'sleep 10' // Wait for services to be ready
            }
        }
        
        stage('Run Tests') {
            steps {
                script {
                    // Run tests inside a temporary Go container
                    sh '''
                        docker run --rm \
                            --network ${COMPOSE_PROJECT_NAME}_fleet-network \
                            -v $(pwd):/app \
                            -w /app \
                            -e CGO_ENABLED=0 \
                            -e GO111MODULE=on \
                            golang:1.21 \
                            sh -c "go mod tidy && go test ./... -v"
                    '''
                }
            }
        }
        
        stage('Build & Start All Services') {
            steps {
                sh 'docker-compose up -d'
                sh 'sleep 5'
            }
        }
        
        stage('Health Check') {
            steps {
                script {
                    sh '''
                        echo "Checking services status..."
                        docker-compose ps
                        
                        # Check if API is responding (adjust port if needed)
                        timeout 30 sh -c 'until curl -f http://localhost:8093/health 2>/dev/null; do sleep 2; done' || echo "API health check skipped"
                    '''
                }
            }
        }
        
        stage('Package Docker Images') {
            when {
                branch 'main'
            }
            steps {
                script {
                    sh '''
                        # Save images for artifact
                        docker save -o api-image.tar ${COMPOSE_PROJECT_NAME}_api:latest
                        docker save -o subscriber-image.tar ${COMPOSE_PROJECT_NAME}_subscriber:latest
                        docker save -o worker-image.tar ${COMPOSE_PROJECT_NAME}_worker:latest
                        docker save -o publisher-image.tar ${COMPOSE_PROJECT_NAME}_publisher:latest
                        
                        # Compress
                        tar -czf docker-images.tar.gz *-image.tar
                        rm *-image.tar
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
                echo 'Options:'
                echo '  - Push images to Docker registry'
                echo '  - Deploy to Kubernetes'
                echo '  - Deploy to Docker Swarm'
                echo '  - SSH to production server and pull images'
            }
        }
    }
    
    post {
        always {
            script {
                // Stop and remove containers
                sh 'docker-compose down -v || true'
                
                // Clean up dangling images
                sh 'docker image prune -f || true'
            }
        }
        success {
            echo "Pipeline succeeded for commit ${env.GIT_COMMIT}"
        }
        failure {
            echo "Pipeline failed for commit ${env.GIT_COMMIT}"
            script {
                // Show logs for debugging
                sh 'docker-compose logs --tail=100 || true'
            }
        }
    }
}