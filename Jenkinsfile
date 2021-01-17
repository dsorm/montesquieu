pipeline {
  agent any
  tools {
          go 'Go 1.15.6'
  }
  environment {
          GO111MODULE = 'on'
          CGO_ENABLED = 0

          PGHOST = "localhost"
          PGPORT = 5005
          PGUSER = test
          PGPASSWORD = test
          PGDATABASE = test
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
              docker.image("postgres").withRun("docker run -d -p ${env.PGPORT}:5432 -e POSTGRES_USER=${env.PGUSER} -e POSTGRES_PASSWORD=${env.PGPASSWORD}")
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