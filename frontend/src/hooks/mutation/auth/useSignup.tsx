import signup from "@/api/auth/signup";
import { useToast } from "@/components/ui/use-toast";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const useSignup = () => {
  const router = useRouter();
  const { toast } = useToast();

  return useMutation({
    mutationFn: signup,
    onSuccess: () => {
      router.push("/signin");
      toast({ description: "회원가입 완료!!" });
    },
  });
};

export default useSignup;
