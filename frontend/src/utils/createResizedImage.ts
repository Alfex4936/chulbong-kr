function changeHeightWidth(
  height: number,
  maxHeight: number,
  width: number,
  maxWidth: number,
  minWidth?: number,
  minHeight?: number
) {
  if (width > maxWidth) {
    height = Math.round((height * maxWidth) / width);
    width = maxWidth;
  }
  if (height > maxHeight) {
    width = Math.round((width * maxHeight) / height);
    height = maxHeight;
  }
  if (minWidth && width < minWidth) {
    height = Math.round((height * minWidth) / width);
    width = minWidth;
  }
  if (minHeight && height < minHeight) {
    width = Math.round((width * minHeight) / height);
    height = minHeight;
  }
  return { height, width };
}

function resizeAndRotateImage(
  image: HTMLImageElement,
  maxWidth: number,
  maxHeight: number,
  minWidth: undefined | number,
  minHeight: undefined | number,
  compressFormat = "jpeg",
  quality = 100,
  rotation = 0
) {
  const qualityDecimal = quality / 100;
  const canvas = document.createElement("canvas");

  let width = image.width;
  let height = image.height;

  const newHeightWidth = changeHeightWidth(
    height,
    maxHeight,
    width,
    maxWidth,
    minWidth,
    minHeight
  );
  if (rotation && (rotation === 90 || rotation === 270)) {
    canvas.width = newHeightWidth.height;
    canvas.height = newHeightWidth.width;
  } else {
    canvas.width = newHeightWidth.width;
    canvas.height = newHeightWidth.height;
  }

  width = newHeightWidth.width;
  height = newHeightWidth.height;

  const ctx = canvas.getContext("2d");

  if (!ctx) {
    return;
  }
  ctx.fillStyle = "rgba(0, 0, 0, 0)";
  ctx.fillRect(0, 0, width, height);

  if (ctx.imageSmoothingEnabled && ctx.imageSmoothingQuality) {
    ctx.imageSmoothingQuality = "high";
  }

  if (rotation) {
    ctx.rotate((rotation * Math.PI) / 180);
    if (rotation === 90) {
      ctx.translate(0, -canvas.width);
    } else if (rotation === 180) {
      ctx.translate(-canvas.width, -canvas.height);
    } else if (rotation === 270) {
      ctx.translate(-canvas.height, 0);
    } else if (rotation === 0 || rotation === 360) {
      ctx.translate(0, 0);
    }
  }
  ctx.drawImage(image, 0, 0, width, height);

  return canvas.toDataURL(`image/${compressFormat}`, qualityDecimal);
}

function b64toByteArrays(b64Data: string) {
  const sliceSize = 512;

  const byteCharacters = atob(
    b64Data.toString().replace(/^data:image\/(png|jpeg|jpg|webp);base64,/, "")
  );
  const byteArrays = [];

  for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
    const slice = byteCharacters.slice(offset, offset + sliceSize);

    const byteNumbers = new Array(slice.length);
    for (let i = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i);
    }

    const byteArray = new Uint8Array(byteNumbers);

    byteArrays.push(byteArray);
  }
  return byteArrays;
}

function b64toBlob(b64Data: string, contentType: string = "image/jpeg") {
  const byteArrays = b64toByteArrays(b64Data);
  const blob = new Blob(byteArrays, {
    lastModified: new Date(), // TODO: check problem
    type: contentType,
  } as BlobPropertyBag);
  return blob;
}

function b64toFile(
  b64Data: string,
  fileName: string,
  contentType: string = "image/jpeg"
) {
  const byteArrays = b64toByteArrays(b64Data);
  const file = new File(byteArrays, fileName, {
    lastModified: new Date(), // TODO: check problem
    type: contentType,
  } as BlobPropertyBag);
  return file;
}

function createResizedImage(
  file: File,
  maxWidth: number,
  maxHeight: number,
  compressFormat: string,
  quality: number,
  rotation: number,
  responseUriFunc: (blob: string | File | Blob) => void,
  outputType = "base64",
  minWidth?: number,
  minHeigh?: number
) {
  const reader = new FileReader();
  if (file) {
    if (file.type && !file.type.includes("image")) {
      throw Error("File Is NOT Image!");
    } else {
      reader.readAsDataURL(file);
      reader.onload = () => {
        const image = new Image();
        image.src = reader.result?.toString() || "";
        image.onload = function () {
          const resizedDataUrl = resizeAndRotateImage(
            image,
            maxWidth,
            maxHeight,
            minWidth,
            minHeigh,
            compressFormat,
            quality,
            rotation
          );

          if (!resizedDataUrl) {
            return; // TODO: improve typings
          }

          const contentType = `image/${compressFormat}`;
          switch (outputType) {
            case "blob":
              const blob = b64toBlob(resizedDataUrl, contentType);
              responseUriFunc(blob);
              break;
            case "base64":
              responseUriFunc(resizedDataUrl);
              break;
            case "file":
              const fileName = file.name;
              const fileNameWithoutFormat = fileName
                .toString()
                .replace(/(png|jpeg|jpg|webp)$/i, "");
              const newFileName = fileNameWithoutFormat.concat(
                compressFormat.toString()
              );
              const newFile = b64toFile(
                resizedDataUrl,
                newFileName,
                contentType
              );
              responseUriFunc(newFile);
              break;
            default:
              responseUriFunc(resizedDataUrl);
          }
        };
      };
      reader.onerror = (error: any) => {
        throw Error(error);
      };
    }
  } else {
    throw Error("File Not Found!");
  }
}

export { createResizedImage };
