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
import useLogin from "@/hooks/mutation/auth/useLogin";
import { zodResolver } from "@hookform/resolvers/zod";
import { isAxiosError } from "axios";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
  email: z.string().email({
    message: "유효한 이메일을 입력해 주세요.",
  }),
  password: z.string().min(8, {
    message: "8자 이상으로 입력해 주세요.",
  }),
});

const SigninForm = () => {
  const { mutateAsync: login, isPending } = useLogin();

  const [errorMessage, setErrorMessage] = useState("");

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await login(values);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          setErrorMessage("유효하지 않은 회원 정보입니다.");
        } else {
          setErrorMessage("잠시 후 다시 시도해 주세요.");
        }
      } else {
        setErrorMessage("잠시 후 다시 시도해 주세요.");
      }
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>이메일</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
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
              <FormMessage>{errorMessage}</FormMessage>
            </FormItem>
          )}
        />
        <Button
          type="submit"
          className="bg-black-light-2 hover:bg-black-light"
          disabled={isPending}
        >
          {isPending ? <LoadingSpinner size="xs" /> : "로그인"}
        </Button>
      </form>
    </Form>
  );
};

export default SigninForm;
