FROM scratch

ADD autograder autograder
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/autograder"]
