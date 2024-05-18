"use client";

import LoadingSpinner from "@/components/atom/LoadingSpinner";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import useSendVerifyCode from "@/hooks/mutation/auth/useSendVerifyCode";
import useSignup from "@/hooks/mutation/auth/useSignup";
import useVerifyCode from "@/hooks/mutation/auth/useVerifyCode";
import isValidEmail from "@/utils/isValidEmail";
import { zodResolver } from "@hookform/resolvers/zod";
import { isAxiosError } from "axios";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import Count from "./Count";
import ErrorMessage from "@/components/atom/ErrorMessage";

const formSchema = z
  .object({
    email: z.string().email({
      message: "유효한 이메일을 입력해 주세요.",
    }),
    username: z
      .string()
      .max(8, {
        message: "8자 이하로 입력해 주세요.",
      })
      .min(2, {
        message: "2자 이상으로 입력해 주세요.",
      }),
    password: z.string().min(8, {
      message: "8자 이상으로 입력해 주세요.",
    }),
    verifyPassword: z.string().min(8, {
      message: "8자 이상으로 입력해 주세요.",
    }),
    code: z.string().length(6, {
      message: "인증 코드는 6자리 숫자여야 합니다.",
    }),
  })
  .refine((data) => data.password === data.verifyPassword, {
    message: "비밀번호가 일치하지 않습니다.",
    path: ["verifyPassword"],
  });

const SignupForm = () => {
  const router = useRouter();

  const { mutateAsync: signup, isPending: signupPending } = useSignup();
  const { mutateAsync: sendCode, isPending: sendCodePending } =
    useSendVerifyCode();
  const { mutateAsync: verifyCode, isPending: verifyPending } = useVerifyCode();

  const [errorMessage, setErrorMessage] = useState("");
  const [emailErrorMessage, setEmailErrorMessage] = useState("");
  const [codeErrorMessage, setCodeErrorMessage] = useState("");

  const [emailBtnText, setEmailBtnText] = useState("인증 요청");
  const [countStart, setCountStart] = useState(false);

  const [isSended, setIsSended] = useState(false);

  const [isVerified, setIsVerified] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
      verifyPassword: "",
      code: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await signup({
        email: values.email,
        password: values.password,
        username: values.username,
      });
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 400) {
          setEmailErrorMessage("이메일 인증이 완료되지 않았습니다.");
          setCodeErrorMessage("이메일 인증이 완료되지 않았습니다.");
        } else {
          setErrorMessage("잠시 후 다시 시도해주세요.");
        }
      } else {
        setErrorMessage("잠시 후 다시 시도해주세요.");
      }
    }
  };

  const handleSendCode = async () => {
    setCountStart(false);
    if (!isValidEmail(form.getValues().email)) {
      setEmailErrorMessage("유효한 이메일을 입력해 주세요.");
      return;
    }
    try {
      await sendCode(form.getValues().email);
      setCountStart(true);
      setEmailBtnText("다시 요청");
      setEmailErrorMessage("");
      setIsSended(true);
      setCountStart(true);
      setIsVerified(false);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 409) {
          setEmailErrorMessage("이미 가입된 이메일입니다.");
        } else {
          setEmailErrorMessage("유효한 이메일을 입력해 주세요.");
        }
      } else {
        setEmailErrorMessage("유효한 이메일을 입력해 주세요.");
      }
    }
  };

  const handleVerify = async () => {
    try {
      await verifyCode({
        code: form.getValues().code,
        email: form.getValues().email,
      });
      setIsVerified(true);
      setCountStart(false);
      setCodeErrorMessage("");
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 400) {
          setCodeErrorMessage("유효하지 않거나 시간이 만료된 코드입니다.");
        } else {
          setCodeErrorMessage("잠시 후 다시 시도해주세요.");
        }
      } else {
        setCodeErrorMessage("유효한 이메일을 입력해 주세요.");
      }
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-3">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>닉네임</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>이메일</FormLabel>
              <div className="relative">
                <FormControl className="">
                  <Input {...field} />
                </FormControl>
                <Button
                  type="button"
                  className="absolute right-0 top-0 text-xs border border-grey border-solid rounded-sm hover:bg-white-tp-dark hover:text-black"
                  onClick={handleSendCode}
                  disabled={sendCodePending}
                >
                  {sendCodePending ? (
                    <LoadingSpinner size="xs" />
                  ) : (
                    <> {emailBtnText}</>
                  )}
                </Button>
              </div>
              <FormMessage>{emailErrorMessage}</FormMessage>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="code"
          render={({ field }) => (
            <FormItem>
              <FormLabel>인증 코드</FormLabel>
              <div className="relative">
                <FormControl className="">
                  <Input {...field} type="text" maxLength={6} />
                </FormControl>
                {isSended && (
                  <Count
                    className="text-xs absolute top-1/2 -translate-y-1/2 right-24"
                    start={countStart}
                    setStart={setCountStart}
                    initTime={300}
                  />
                )}

                <Button
                  type="button"
                  className="absolute right-0 top-0 text-xs border border-grey border-solid rounded-sm hover:bg-white-tp-dark hover:text-black"
                  onClick={handleVerify}
                  disabled={isVerified}
                >
                  {verifyPending ? (
                    <LoadingSpinner size="xs" />
                  ) : isVerified ? (
                    "✅"
                  ) : (
                    "인증 확인"
                  )}
                </Button>
              </div>
              <FormMessage>{codeErrorMessage}</FormMessage>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>비밀번호</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="verifyPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>비밀번호 확인</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div>
          <Button
            type="submit"
            className="bg-black-light-2 mr-3 hover:bg-black-light"
            disabled={signupPending}
          >
            {signupPending ? <LoadingSpinner size="xs" /> : "회원가입"}
          </Button>
          <Button
            type="button"
            className="bg-grey-dark text-black hover:bg-grey-light"
            onClick={() => router.push("/signin")}
          >
            취소
          </Button>
        </div>
        <ErrorMessage text={errorMessage} />
      </form>
    </Form>
  );
};

export default SignupForm;
