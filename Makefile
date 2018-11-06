all: autoscale
	@source setgopath.sh && go build autoscale.go

linux: autoscale
	@source setgopath.sh && GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build autoscale.go

run: all
	@source aws-env.sh && ./autoscale

docker:
	@docker build -t russokj/aws-asg:latest .
	@docker push russokj/aws-asg

clean:
	@rm autoscale

autoscale: autoscale.go


