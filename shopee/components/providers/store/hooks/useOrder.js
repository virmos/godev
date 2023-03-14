import useSWR from "swr"
import axios from "axios";

export const handler = () => (orderId) => {
    const swrRes = useSWR(`/api/order`,
        async () => {
            const { data } = await axios.get(`/api/orders/${orderId}`);
            return data
        }
    )
  
    return swrRes
}