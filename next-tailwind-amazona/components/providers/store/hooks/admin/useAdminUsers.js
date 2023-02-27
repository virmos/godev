import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/admin/users`,
        async () => {
            const { data } = await axios.get(`/api/admin/users`);
            return data
        }
    )
  
    return swrRes
}