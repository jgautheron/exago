FROM scratch
ADD exago /
CMD ["-h"]
ENTRYPOINT ["/exago"]
EXPOSE 3000
