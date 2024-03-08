import styled from "@emotion/styled";

export const ImageUploadContainer = styled.div`
  margin-bottom: 2rem;

  & > input {
    display: none;
  }
`;

export const ImageBox = styled.div`
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

  width: 170px;
  height: 170px;

  cursor: pointer;
`;
