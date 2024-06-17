import { getUser } from "@/actions/auth";
import LogoutButton from "./logout-button";

const ProtectedPage = async () => {
  const user = await getUser();

  return (
    <div>
      <h1>This is a protected page</h1>
      <p>{JSON.stringify(user)}</p>
      <LogoutButton />
    </div>
  );
};

export default ProtectedPage;
