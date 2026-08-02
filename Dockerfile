FROM golang:1.21-bookworm

LABEL maintainer="QuotientBot Team"

# Install runtime dependencies and build tools
RUN apt-get update -qq && apt-get install -y -qq \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr \
    ca-certificates \
    curl \
    git \
 && rm -rf /var/lib/apt/lists/*

# Set up Tessdata directory
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
RUN mkdir -p ${TESSDATA_PREFIX}

# Download high-speed tessdata_fast models (sub-200ms recognition speed)
RUN curl -L -o ${TESSDATA_PREFIX}eng.traineddata https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata \
 && curl -L -o ${TESSDATA_PREFIX}osd.traineddata https://github.com/tesseract-ocr/tessdata_fast/raw/main/osd.traineddata

# Configure OpenMP to prevent CPU busy-spinning / thread hogging
ENV OMP_NUM_THREADS=1
ENV OMP_THREAD_LIMIT=1
ENV OMP_WAIT_POLICY=PASSIVE
ENV OMP_DYNAMIC=FALSE

WORKDIR /app

COPY . .

# Regenerate module checksums and build binary
RUN go mod tidy && go build -v -o /app/ocr_v2 .

EXPOSE 8080
CMD ["/app/ocr_v2"]
