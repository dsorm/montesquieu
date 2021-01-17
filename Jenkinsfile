pipeline {
  agent any
  tools {
          go 'Go 1.15.6'
  }
  environment {
          GO111MODULE = 'on'
          CGO_ENABLED = 0
  }
  stages {
    stage('Run') {
      parallel {
        stage('Go Build') {
          agent any
          steps {
            echo 'Building...'
            catchError(buildResult: 'Failure') {
              sh 'go build -o montesquieu .'
            }

          }
        }

        stage('Go Test') {
          agent any
          steps {
            echo 'Testing...'
            warnError(message: 'Unstable', catchInterruptions: true) {
              sh 'go test ./...'
            }

          }
        }

      }
    }

    stage('Archive artifacts') {
      agent any
      steps {
        archiveArtifacts(artifacts: 'montesquieu')
        echo 'Artifacts archived'
      }
    }

  }
}