function add
    http :8000/event type=create name=test kind=cr
    http :8000/event type=create name=test kind=crd
    http :8000/event type=create name=test1 kind=cr
end

add
