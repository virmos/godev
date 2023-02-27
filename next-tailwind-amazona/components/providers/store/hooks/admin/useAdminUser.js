import useSWR from "swr"
import axios from "axios";

export const handler = () => (userId) => {
    const swrRes = useSWR(`/api/admin/users/${userId}`,
        async () => {
            const { data } = await axios.get(`/api/admin/users/${userId}`);
            return data
        }
    )
  
    return swrRes
}