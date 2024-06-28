import {
  Center,
  Heading,
} from '@chakra-ui/react'
import CustomerManagement from '../../components/CustomerManagement'

export default function Home() {
  return (
    <>
      <Center>
        <Heading margin="20px">
          Customer Management
        </Heading>
      </Center>
      <CustomerManagement />
    </>
  )
}
