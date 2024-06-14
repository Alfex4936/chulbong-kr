import "@testing-library/jest-dom";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import UploadImage from "@/app/(navbar)/pullup/register/_components/UploadImage";

class MockFileReader {
  onloadend: () => void;
  readAsDataURL: (blob: Blob) => void;
  result: string | ArrayBuffer | null;

  constructor() {
    this.onloadend = jest.fn();
    this.readAsDataURL = jest.fn(() => {
      this.result = "data:image/png;base64,...";
      this.onloadend();
    });
    this.result = null;
  }

  static EMPTY = 0;
  static LOADING = 1;
  static DONE = 2;
}

window.FileReader = MockFileReader as any;

Object.defineProperty(globalThis.Image.prototype, "src", {
  set() {
    setTimeout(() => {
      if (typeof this.onload === "function") {
        this.onload();
      }
    }, 0);
  },
});

jest.mock("uuid", () => ({
  v4: () => "12345",
}));

jest.mock("@/utils/resizeFile", () => ({
  __esModule: true,
  default: jest.fn().mockImplementation((file) =>
    Promise.resolve(
      new File([file], file.name, {
        type: file.type,
      })
    )
  ),
}));

describe("UploadImage Component", () => {
  beforeEach(() => {
    render(<UploadImage />);
  });
  it("should display an error message for an unsupported file format.", async () => {
    const input = screen.getByTestId("file-input");

    const file = new File(["(๑•̀ㅂ•́)و✧"], "unsupportedFormat.pdf", {
      type: "image/pdf",
    });

    fireEvent.change(input, { target: { files: [file] } });

    const errorMessage = screen.getByTestId("file-error");
    expect(errorMessage).toHaveTextContent(
      "지원되지 않은 이미지 형식입니다. JPEG, PNG, webp형식의 이미지를 업로드해주세요."
    );
  });

  it("should display 2 buttons with aria-label='삭제' when 2 images are uploaded.", async () => {
    const input = screen.getByTestId("file-input");

    const files = [
      new File(["(๑•̀ㅂ•́)و✧1"], "1.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧2"], "2.png", { type: "image/png" }),
    ];

    for (const file of files) {
      await waitFor(() => {
        fireEvent.change(input, { target: { files: [file] } });
      });
    }

    const deleteButtons = screen.getAllByText("삭제");
    expect(deleteButtons.length).toBe(2);
  });

  it("should display an error message if more than 5 images are uploaded.", async () => {
    const input = screen.getByTestId("file-input");

    const files = [
      new File(["(๑•̀ㅂ•́)و✧1"], "1.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧2"], "2.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧3"], "3.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧4"], "4.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧5"], "5.png", { type: "image/png" }),
      new File(["(๑•̀ㅂ•́)و✧6"], "6.png", { type: "image/png" }),
    ];

    for (const file of files) {
      await waitFor(() => {
        fireEvent.change(input, { target: { files: [file] } });
      });
    }

    const errorMessage = screen.getByTestId("file-error");
    expect(errorMessage).toHaveTextContent("최대 5개 까지 등록 가능합니다!");
  });
});
