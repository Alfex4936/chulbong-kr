import styled from "@emotion/styled";

export const ImageUploadContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;

  margin-bottom: 1rem;

  & > div > input {
    display: none;
  }
`;

export const ImageBox = styled.div`
  position: relative;
  border-radius: 0.5rem;

  display: flex;
  justify-content: center;
  align-items: center;

  background-image: ${({ img }: { img: string | null }) => {
    return img && `url(${img})`;
  }};

  object-fit: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;

  margin: auto;

  width: 40px;
  height: 40px;

  cursor: pointer;
`;

export const ImageViewContainer = styled.div`
  margin-top: 0.5rem;
  display: flex;
  justify-content: center;
`;

export const ImageView = styled.div`
  width: 55px;
  height: 55px;

  margin: 0.3rem;

  background-image: ${({ img }: { img: string | null }) => {
    return img && `url(${img})`;
  }};
  object-fit: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;

  cursor: pointer;

  border: 1px solid #ccc;
  border-radius: 0.4rem;
`;

export const ErrorBox = styled.div`
  text-align: center;

  font-size: 0.7rem;

  color: red;
`;
