import * as Styles from "./SignupFormSkeleton.style";

const SignupFormSkeleton = () => {
  return (
    <div>
      <Styles.TitleSkeleton />
      <Styles.InputSkeleton />
      <Styles.InputSkeleton />
      <Styles.InputSkeleton />
      <Styles.InputSkeleton />
      <Styles.ButtonSkeleton />
    </div>
  );
};

export default SignupFormSkeleton;
