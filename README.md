# On

This is a command line tool helps you avoid keep typing the same prefix when you are operating the very similar commands

For example:
```
$ on kubectl
on(kubectl)> get po
NAME                                READY   STATUS    RESTARTS   AGE
alpine-deploy-7b496f9bb-rdt6q       1/1     Running   2          28h
alpine-deploy-7b496f9bb-stcqj       1/1     Running   2          28h
api-gateway-75d7c4b897-cb2gp        1/1     Running   1          12d
curl-74dbc9bc95-7qcmn               1/1     Running   1          12d
curl-74dbc9bc95-9j5dx               1/1     Running   1          12d
hello-deployment-6d474d7946-96dtf   1/1     Running   0          28h
hello-deployment-6d474d7946-9dshj   1/1     Running   0          28h
hello-deployment-6d474d7946-j52dp   1/1     Running   0          28h
world-deployment-6886f8cf9f-c4wlm   1/1     Running   0          28h
world-deployment-6886f8cf9f-dz4j5   1/1     Running   0          28h
world-deployment-6886f8cf9f-zvfx9   1/1     Running   0          28h

on(kubectl)> get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
hello-svc    ClusterIP   10.0.154.37    <none>        80/TCP    28h
kubernetes   ClusterIP   10.0.0.1       <none>        443/TCP   12d
world-svc    ClusterIP   10.0.106.133   <none>        80/TCP    28h
```

### Key Binds

- `<C-a>`: control + A, type this then write done new contexts and `<Enter>`, you would see the context be apply on to old context

    For example:
    ```
    on(kubectl)>
    # <C-a>get<Enter>
    on(kubectl get)>
    ```
- `<C-c>`: control + C, type this would pop out the last element in command context

    For example:
    ```
    on(kubectl get)>
    # <C-c>
    on(kubectl)>
    ```

At here you can see we avoid typing `kubectl` again and again
