task :default => [:all]
task :build => [:docker_build]
task :all => [ :docker_build]

task docker_build: [:go_build ] do
    STDOUT.puts "What should the Autograder image version tag be?"
    input = STDIN.gets.strip

    if input == ''
        Rake::Task["docker_build"].reenable
        Rake::Task["docker_build"].invoke
    else
        sh "docker build -t gcr.io/alligrader-15/autograder:#{input} ."
        sh "gcloud docker push gcr.io/alligrader-15/autograder:#{input}"
    end
end

task go_build: [ ] do
    sh "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ."
end

task deploy: [] do
    sh "kubectl run autograder --image=gcr.io/alligrader-15/autograder:v0.0.2 --port=80"
end
