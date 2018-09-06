# Run test inside docker
FROM neojrotary/gcb-bridge/test-base
COPY . .
CMD ["bash", "./tests.sh"]