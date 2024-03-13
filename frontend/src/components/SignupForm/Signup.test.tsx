import "@testing-library/jest-dom";
import { fireEvent, render, screen } from "@testing-library/react";
import axios from "axios";
import SignupForm from "./SignupForm";
import { act } from "react-dom/test-utils";

jest.mock("axios");
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe("signup test", () => {
  beforeEach(() => {
    render(<SignupForm />);
  });
  it("Should display an error message for an invalid email format.", async () => {
    // given

    // when
    const emailInput = screen.getByTestId("email");
    const signupButton = screen.getByTestId("signup-button");

    fireEvent.change(emailInput, { target: { value: "notemailvalue" } });
    fireEvent.click(signupButton);

    // then
    const errorMessage = await screen.findByTestId("email-error");
    expect(errorMessage).toHaveTextContent("이메일 형식이 아닙니다.");
  });

  it("Should show the verification code input when a valid email is entered.", async () => {
    // given
    const responseData = "Verification code sent successfully.";
    mockedAxios.post.mockResolvedValue({ msg: responseData });

    // when
    const emailInput = screen.getByTestId("email");
    const sendCodeButton = screen.getByText("인증 요청");

    await act(async () => {
      fireEvent.change(emailInput, { target: { value: "asd@asd.com" } });
      fireEvent.click(sendCodeButton);
    });
    // then
    const codeInput = screen.queryByText("인증번호");
    expect(codeInput).toBeInTheDocument();
  });

  it("Should display an error message for an invalid password format.", async () => {
    // given

    // when
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");

    fireEvent.change(passwordInput, { target: { value: "password" } });
    fireEvent.change(verifyPasswordInput, { target: { value: "password" } });
    fireEvent.click(signupButton);
    // then
    const errorMessage = await screen.findByTestId("password-error");
    expect(errorMessage).toHaveTextContent(
      "특수문자 포함 8 ~ 20자 사이로 입력해 주세요."
    );
  });

  it("Should display an error message if the password confirmation does not match.", async () => {
    // given
    const responseData = "Signup successfully.";
    mockedAxios.post.mockResolvedValue({ msg: responseData });

    // when
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");

    fireEvent.change(passwordInput, { target: { value: "password1234!" } });
    fireEvent.change(verifyPasswordInput, {
      target: { value: "password4321!" },
    });
    fireEvent.click(signupButton);
    // then
    const errorMessage = await screen.findByTestId("verify-password-error");
    expect(errorMessage).toHaveTextContent("비밀번호를 확인해 주세요.");
  });

  it("Should not display the signup form when registration is successful.", async () => {
    // given
    const responseData = "Signup successfully.";
    mockedAxios.post.mockResolvedValue({ msg: responseData });

    // when
    const nameInput = screen.getByTestId("name");
    const emailInput = screen.getByTestId("email");
    const passwordInput = screen.getByTestId("password");
    const verifyPasswordInput = screen.getByTestId("verify-password");

    const signupButton = screen.getByTestId("signup-button");
    const sendCodeButton = screen.getByText("인증 요청");

    fireEvent.change(nameInput, { target: { value: "testName" } });

    await act(async () => {
      fireEvent.change(emailInput, { target: { value: "asd@asd.com" } });
      fireEvent.click(sendCodeButton);
    });

    await act(async () => {
      const codeInput = screen.getByTestId("code");
      fireEvent.change(codeInput, { target: { value: "123456" } });

      const confirmButton = screen.getByText("인증 확인");
      fireEvent.click(confirmButton);
    });

    await act(async () => {
      fireEvent.change(passwordInput, { target: { value: "password4321!" } });
      fireEvent.change(verifyPasswordInput, {
        target: { value: "password4321!" },
      });

      fireEvent.click(signupButton);
    });

    // then
    const signupText = screen.getByTestId("signup-success");
    expect(signupText).toHaveTextContent("회원 가입 완료!");
  });
});
