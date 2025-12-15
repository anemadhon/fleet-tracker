pipeline {
    agent any
    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '20'))
    }
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Set up Go') {
            steps {
                sh 'go version'
                sh 'go env'
            }
        }
        stage('Download deps') {
            steps {
                sh 'go mod tidy'
            }
        }
        stage('Unit tests') {
            steps {
                sh 'go test ./... -v'
            }
        }
        stage('Build API binary') {
            steps {
                sh 'go build -o bin/api ./cmd/api'
            }
        }
        stage('Package artifact') {
            when {
                expression { fileExists('bin/api') }
            }
            steps {
                sh 'tar -czf api-artifact.tar.gz bin/api'
                archiveArtifacts artifacts: 'api-artifact.tar.gz', fingerprint: true
            }
        }
        stage('Deploy (placeholder)') {
            when {
                branch 'main'
            }
            steps {
                echo 'Deploy your Go API here (e.g. Docker build & push, ssh, k8s, etc.)'
            }
        }
    }
    post {
        success {
            echo "Go API pipeline succeeded for commit ${env.GIT_COMMIT}"
        }
        failure {
            echo "Go API pipeline failed for commit ${env.GIT_COMMIT}"
        }
    }
}
