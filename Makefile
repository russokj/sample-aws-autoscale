all: autoscale
	@source setgopath.sh && go build autoscale.go

run: all
	@source aws-env.sh && ./autoscale

clean:
	@rm autoscale

autoscale: autoscale.go


