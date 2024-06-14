import SignupForm from "@/app/(non-navbar)/signup/_components/signup-form";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@testing-library/jest-dom";
import { fireEvent, render, screen } from "@testing-library/react";
import axios, { type AxiosError } from "axios";
import { act } from "react-dom/test-utils";

jest.mock("axios", () => {
  const mockAxiosInstance = {
    post: jest.fn(),
    get: jest.fn(),
    interceptors: {
      request: { use: jest.fn(), eject: jest.fn() },
      response: { use: jest.fn(), eject: jest.fn() },
    },
  };

  return {
    create: jest.fn(() => mockAxiosInstance),
    isAxiosError: (error: any): error is AxiosError => true,
  };
});

const mockedAxios = axios.create() as jest.Mocked<typeof axios>;

const queryClient = new QueryClient();

describe("회원가입 테스트", () => {
  beforeEach(() => {
    render(
      <QueryClientProvider client={queryClient}>
        <SignupForm />
      </QueryClientProvider>
    );
  });

  it("유효하지 않은 이메일 에러", async () => {
    const emailInput = screen.getByTestId("email");
    const signupButton = screen.getByTestId("signup-button");

    fireEvent.change(emailInput, { target: { value: "notemailvalue" } });
    fireEvent.click(signupButton);

    const errorMessage = await screen.findByTestId("email-error");
    expect(errorMessage).toHaveTextContent("유효한 이메일을 입력해 주세요.");
  });

  it("이메일 보내기 성공 후 다시 요청 활성화", async () => {
    const responseData = { msg: "success" };
    mockedAxios.post.mockResolvedValueOnce({ data: responseData });

    const emailInput = screen.getByTestId("email");
    const sendCodeButton = screen.getByTestId("send-email-btn");

    await act(async () => {
      fireEvent.change(emailInput, { target: { value: "asd@asd.com" } });
      fireEvent.click(sendCodeButton);
    });

    expect(sendCodeButton).toHaveTextContent("다시 요청");
  });

  it("비밀번호 일치 여부", async () => {
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");

    await act(async () => {
      fireEvent.change(passwordInput, { target: { value: "password" } });
      fireEvent.change(verifyPasswordInput, { target: { value: "passwordd" } });
      fireEvent.click(signupButton);
    });

    const errorMessage = await screen.findByTestId("verify-password-error");
    expect(errorMessage).toHaveTextContent("비밀번호가 일치하지 않습니다.");
  });

  it("비밀번호 8자 이상", async () => {
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");

    await act(async () => {
      fireEvent.change(passwordInput, { target: { value: "passwo" } });
      fireEvent.change(verifyPasswordInput, { target: { value: "passwo" } });
      fireEvent.click(signupButton);
    });

    const errorMessage = await screen.findByTestId("password-error");
    expect(errorMessage).toHaveTextContent("8자 이상으로 입력해 주세요.");
  });

  it("이메일 인증 완료 여부", async () => {
    mockedAxios.post.mockRejectedValueOnce({
      response: { status: 400 },
    });

    const nameInput = screen.getByTestId("username");
    const emailInput = screen.getByTestId("email");
    const codeInput = screen.getByTestId("code");
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");

    await act(async () => {
      fireEvent.change(nameInput, { target: { value: "name" } });
      fireEvent.change(emailInput, { target: { value: "asd@asd.com" } });
      fireEvent.change(codeInput, { target: { value: "123456" } });
      fireEvent.change(passwordInput, { target: { value: "password" } });
      fireEvent.change(verifyPasswordInput, { target: { value: "password" } });

      fireEvent.click(signupButton);
    });

    const errorMessage = await screen.findByTestId("email-error");
    expect(errorMessage).toHaveTextContent(
      "이메일 인증이 완료되지 않았습니다."
    );
  });
});
