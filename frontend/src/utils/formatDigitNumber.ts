const DIGIT = 2;

const formatDigitNumber = (number: number) => {
  return `${number.toString().padStart(DIGIT, "0")}`;
};

export default formatDigitNumber;
