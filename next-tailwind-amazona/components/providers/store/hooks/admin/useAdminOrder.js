import useSWR from "swr"
import axios from "axios";

export const handler = () => (orderId) => {
    const swrRes = useSWR(`/api/admin/orders/${orderId}`,
        async () => {
            const { data } = await axios.get(`/api/admin/orders/${orderId}`);
            return data
        }
    )
  
    return swrRes
}