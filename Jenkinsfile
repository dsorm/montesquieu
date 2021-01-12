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
            sh 'go build -o montesquieu .'
            catchError(buildResult: 'Failure')
          }
        }

        stage('Go Test') {
          agent any
          steps {
            echo 'Testing...'
            sh 'go test ./...'
            warnError(message: 'Unstable', catchInterruptions: true)
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