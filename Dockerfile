FROM scratch
ADD exago /
CMD ["-h"]
ENTRYPOINT ["/exago"]