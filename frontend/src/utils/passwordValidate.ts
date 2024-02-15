const passwordValidate = (password: string) => {
  const regex =
    /^(?=.*[a-zA-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/;

  return regex.test(password);
};

export default passwordValidate;
