EXEC = readline

all:bin

bin:
	@go build -o $(EXEC)

install:$(EXEC)
	install -d $(BIN) $(COMP_DIR)
	install $(EXEC) $(BIN)

clean:
	@rm -rfv $(EXEC)

archive:
	@echo "archive to $(TAR)"
	@git archive master --prefix="$(EXEC)-$(VER)/" --format tar.gz -o $(TAR)

test:
	@go test
