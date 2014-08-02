FROM       scratch
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
ADD        contributors contributors
ENV        PORT 80
EXPOSE     80
ENTRYPOINT ["/contributors"]
