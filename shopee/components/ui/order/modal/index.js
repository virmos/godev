import { Modal, Button } from "@components/ui/common";
import { useEffect, useState } from "react";

const defaultProduct = {
    name: '',
    slug: '',
    price: 0,
}

const _createFormState = (isDisabled = false, message =    "") => ({isDisabled, message})

const createFormState = ({name, slug, price}) => {
    if (!price || Number(price) <= 0) {
        return _createFormState(true, "Price is not valid.")
    }
    if (!name || name.length < 3) {
        return _createFormState(true, "Name is not valid.")
    }
    if (!slug) {
        return _createFormState(true, "Slug is not valid.")
    }

    return _createFormState()
}

export default function ProductModal({trigger, onClose, onSubmit}) {
    const [isOpen, setIsOpen] = useState(false)
    const [product, setProduct] = useState(defaultProduct)

    useEffect(() => {
        setIsOpen(true)
        setProduct({
            ...defaultProduct,
        })
    }, [trigger])

    const closeModal = () => {
        setIsOpen(false)
        setProduct(defaultProduct)
        onClose()
    }

    const formState = createFormState(product)

    return (
        <Modal isOpen={isOpen}>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                    <div className="sm:flex sm:items-start">
                        <div className="mt-3 sm:mt-0 sm:ml-4 sm:text-left">
                            <h3 className="mb-7 text-lg font-bold leading-6 text-gray-900" id="modal-title">
                                Fill in product form
                            </h3>
                            <div className="mt-1 relative rounded-md">
                                <label className="mb-2 font-bold"> Name</label>
                                <input
                                    value={product.name}
                                    onChange={({target: {value}}) => {
                                        setProduct({
                                            ...product,
                                            name: value
                                        })
                                    }}
                                    type="text"
                                    name="name"
                                    id="name"
                                    className="disabled:opacity-50 w-80 mb-1 focus:ring-indigo-500 shadow-md focus:bproduct-indigo-500 block pl-7 p-4 sm:text-sm bproduct-gray-300 rounded-md"
                                    placeholder="sample name"
                                />
                                <p className="text-xs text-gray-700 mt-1">
                                    Name must have at least 3 characters.
                                </p>
                            </div>
                            <div className="mt-1 relative rounded-md">
                                <label className="mb-2 font-bold"> Slug</label>
                                <input
                                    value={product.slug}
                                    onChange={({target: {value}}) => {
                                        setProduct({
                                            ...product,
                                            slug: value.trim()
                                        })
                                    }}
                                    type="text"
                                    name="slug"
                                    id="slug"
                                    className="disabled:opacity-50 w-80 mb-1 focus:ring-indigo-500 shadow-md focus:bproduct-indigo-500 block pl-7 p-4 sm:text-sm bproduct-gray-300 rounded-md"
                                    placeholder="sample-name-1"
                                />
                                <p className="text-xs text-gray-700 mt-1">
                                It&apos;s important to choose a good slug, because it will be used in the URL of your product.
                                </p>
                            </div>
                            <div className="mt-1 relative rounded-md">
                                <label className="mb-2 font-bold"> Price</label>
                                <input
                                    value={product.price}
                                    onChange={({target: {value}}) => {
                                        if (isNaN(value)) { return; }
                                        setProduct({
                                            ...product,
                                            price: value.trim()
                                        })
                                    }}
                                    type="text"
                                    name="price"
                                    id="price"
                                    className="disabled:opacity-50 w-80 mb-1 focus:ring-indigo-500 shadow-md focus:bproduct-indigo-500 block pl-7 p-4 sm:text-sm bproduct-gray-300 rounded-md"
                                    placeholder="0"
                                />
                                <p className="text-xs text-gray-700 mt-1">
                                    Price must be greater than 0.
                                </p>
                            </div>
                            { formState.message &&
                                <div className="p-4 my-3 text-yellow-700 bg-yellow-200 rounded-lg text-sm">
                                    { formState.message }
                                </div>
                            }
                        </div>
                    </div>
                </div>
                <div className="bg-gray-50 px-4 py-3 sm:px-6 flex">
                    <Button
                        disabled={formState.isDisabled}
                        onClick={() => {
                            onSubmit(product)
                    }}>
                        Submit
                    </Button>
                    <Button
                        onClick={closeModal}
                        variant="red">
                        Cancel
                    </Button>
                </div>
            </div>
        </Modal>
    )
}
