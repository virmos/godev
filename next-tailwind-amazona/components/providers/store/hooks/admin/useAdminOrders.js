import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/admin/orders`,
        async () => {
            const { data } = await axios.get(`/api/admin/orders`);
            return data
        }
    )
  
    return swrRes
}