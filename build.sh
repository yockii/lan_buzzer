#!/bin/bash

echo "Building frontend..."
cd frontend
npm run build
cd ..

echo "Copying frontend to backend..."
mkdir -p backend/embed/dist
rm -rf backend/embed/dist/*
cp -r frontend/dist/* backend/embed/dist/

echo "Building backend..."
cd backend
go build -o ../build/lan-buzzer.exe .
cd ..

echo "Build complete: build/lan-buzzer.exe"
