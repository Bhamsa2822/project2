import axios from "axios"
import { Customer } from "../components/customer"

export function fetchAllCustomers(): Promise<Customer[]> {
    return axios.get("/api/customers")
        .then((response) => {
            return response.data
        })
}

export function createCustomer(data: Customer): Promise<string> {
    return axios.post("/api/customers", data)
        .then((res) => { return res.data })
}

export function updateCustomer(data: Customer): Promise<string> {
    return axios.put("/api/customers", data)
        .then((res) => {
            return res.data
        })
}

export function deleteCustomer(id: string): Promise<void> {
    return axios.delete(`/api/customers/${id}`)
}