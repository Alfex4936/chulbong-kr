"use client";

import ErrorMessage from "@/components/atom/ErrorMessage";
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
import useResetPassword from "@/hooks/mutation/auth/useResetPassword";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
  password: z.string().min(8, {
    message: "8자 이상으로 입력해 주세요.",
  }),
});

const ResetPasswordClient = () => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      password: "",
    },
  });

  const { mutate: reset, isPending, isError } = useResetPassword();
  const searchParams = useSearchParams();

  const token = searchParams.get("token");

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    reset({ token: token as string, password: values.password });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>비밀번호 입력</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              {isError ? (
                <ErrorMessage text="인증이 만료됐습니다." />
              ) : (
                <FormMessage />
              )}
            </FormItem>
          )}
        />
        <Button
          type="submit"
          className="disabled:bg-black-light"
          disabled={isPending}
        >
          {isPending ? <LoadingSpinner color="white" size="xs" /> : "변경"}
        </Button>
      </form>
    </Form>
  );
};

export default ResetPasswordClient;
