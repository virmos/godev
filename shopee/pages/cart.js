import Link from "next/link";
import React from "react";
import { BaseLayout } from "@components/ui/layout";
import { useStore } from "@components/providers";
import dynamic from "next/dynamic";
import { toast } from "react-toastify";
import { getProduct } from "@components/api";
import { Button, Message } from "@components/ui/common";
import { CartReport } from "@components/ui/cart";
import { CartTable } from "@components/ui/cart";

function CartScreen() {
    const { state, dispatch } = useStore();
    const {
        cart: { cartItems },
    } = state;
    const removeItemHandler = (item) => {
        dispatch({ type: "CART_REMOVE_ITEM", payload: item });
    };
    const updateCartHandler = async (item, qty) => {
        const quantity = Number(qty);
        const { data } = await getProduct(item._id);
        if (data.countInStock < quantity) {
            return toast.error("Sorry. Product is out of stock");
        }
        dispatch({ type: "CART_ADD_ITEM", payload: { ...item, quantity } });
    };
    return (
        <BaseLayout>
            <h1 className="mb-4 text-xl">Shopping Cart</h1>
            {cartItems.length === 0 ? (
                <Message type="warning">
                    Cart is empty.
                    <Link href="/">
                        <a>
                            <Button size="sm" variant="purple" className="ml-2">
                                Go Shopping
                            </Button>
                        </a>
                    </Link>
                </Message>
            ) : (
                <div className="grid md:grid-cols-4 md:gap-5">
                    <div className="overflow-x-auto md:col-span-3">
                        <CartTable
                            cartItems={cartItems}
                            updateCartHandler={updateCartHandler}
                            removeItemHandler={removeItemHandler}
                        />
                    </div>
                    <CartReport cartItems={cartItems} />
                </div>
            )}
        </BaseLayout>
    );
}

export default dynamic(() => Promise.resolve(CartScreen), { ssr: false });
