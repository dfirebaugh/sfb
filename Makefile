CC = gcc
CFLAGS = -Wall -Wextra -O2
LDFLAGS = -lSDL2 -mavx

TARGET_LIB = libsfb.so
TARGET_TEST = simd_fb_test

SRC = sfb.c simd_draw.c
TEST_SRC = main.c

all: $(TARGET_LIB)
	
test: $(TARGET_TEST)
	
$(TARGET_LIB): $(SRC)
	$(CC) -shared -o $@ -fPIC $(SRC) $(LDFLAGS)

$(TARGET_TEST): $(SRC) $(TEST_SRC)
	$(CC) -o $@ $(TEST_SRC) $(SRC) $(LDFLAGS)

clean:
	rm -f $(TARGET_LIB) $(TARGET_TEST)

.PHONY: all clean
