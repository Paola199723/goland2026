#!/bin/bash

# Script para iniciar el sistema Goland 2026 limpiamente

echo "🧹 PASO 1: Limpiando procesos antiguos..."
pkill -9 api 2>/dev/null
pkill -9 node 2>/dev/null
sleep 2

echo "🧹 PASO 2: Limpiando base de datos..."
psql -h localhost -p 26257 -U root -d market_data -c "DELETE FROM challenges; DELETE FROM page_cursors; DELETE FROM auth_tokens;" 2>/dev/null
echo "✅ Base de datos limpiada"

echo ""
echo "⚙️ PASO 3: Configurando .env..."
cat > /Users/paolacasadiegos/Downloads/truora/goland2026/.env << 'EOF'
# KarenAI API Token - Usar placeholder para mock data
# Si tienes token real, reemplaza aquí
# KarenAI API Token - Usar placeholder para mock data
# Si tienes token real, reemplaza aquí
#TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6NSwiZW1haWwiOiJjYXNhZGllZ29zdmFjYUBnbWFpbC5jb20iLCJleHAiOjE3NzAxNTE4MDEsImlkIjoiIiwicGFzc3dvcmQiOiJ0IFx0RlJPTSBcdCB1c2VycyBcdCBXSEVSRVx0IHVzZXJuYW1lXHQgaXMgXHQgbm90IFx0IG51bGwgXHQgQU5EIFx0IHBhc3N3b3JkIFx0IGlzIFx0IG5vdCBcdCBudWxsIFx0IFVOSU9OIFx0IFNFTEVDVCBcdCB1c2VybmFtZSwgcGFzc3dvcmQgXHQgYXMgXHQgIHQifQ.RcOFHsDjE-CjGRW4_WXvMH20IFYdA1-y7kdB-eqP4R0


# Database Configuration
DATABASE_URL=postgresql://root@localhost:26257/market_data?sslmode=disable

# JWT Secret
#WT_SECRET=your_jwt_secret_here
# KarenAI API Token
TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6NSwiZW1haWwiOiJjYXNhZGllZ29zdmFjYUBnbWFpbC5jb20iLCJleHAiOjE3NzAxNTE4MDEsImlkIjoiIiwicGFzc3dvcmQiOiJ0IFx0RlJPTSBcdCB1c2VycyBcdCBXSEVSRVx0IHVzZXJuYW1lXHQgaXMgXHQgbm90IFx0IG51bGwgXHQgQU5EIFx0IHBhc3N3b3JkIFx0IGlzIFx0IG5vdCBcdCBudWxsIFx0IFVOSU9OIFx0IFNFTEVDVCBcdCB1c2VybmFtZSwgcGFzc3dvcmQgXHQgYXMgXHQgIHQifQ.RcOFHsDjE-CjGRW4_WXvMH20IFYdA1-y7kdB-eqP4R0

# JWT Token expiration time (in hours)
JWT_EXPIRATION=24

# JWT Secret key
JWT_SECRET=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6NSwiZW1haWwiOiJjYXNhZGllZ29zdmFjYUBnbWFpbC5jb20iLCJleHAiOjE3NzAxNTE4MDEsImlkIjoiIiwicGFzc3dvcmQiOiJ0IFx0RlJPTSBcdCB1c2VycyBcdCBXSEVSRVx0IHVzZXJuYW1lXHQgaXMgXHQgbm90IFx0IG51bGwgXHQgQU5EIFx0IHBhc3N3b3JkIFx0IGlzIFx0IG5vdCBcdCBudWxsIFx0IFVOSU9OIFx0IFNFTEVDVCBcdCB1c2VybmFtZSwgcGFzc3dvcmQgXHQgYXMgXHQgIHQifQ.RcOFHsDjE-CjGRW4_WXvMH20IFYdA1-y7kdB-eqP4R0
s
EOF
echo "✅ .env configurado con mock data"

echo ""
echo "🔨 PASO 4: Compilando backend..."
cd /Users/paolacasadiegos/Downloads/truora/goland2026/backendgoland2026
go build -o api ./cmd/api
if [ $? -ne 0 ]; then
    echo "❌ Error compilando backend"
    exit 1
fi
echo "✅ Backend compilado"

echo ""
echo "🚀 PASO 5: Iniciando backend en puerto 8081..."
./api > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
sleep 3

# Verificar que el backend está corriendo
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "❌ Backend no se inició correctamente"
    echo "Ver logs: tail /tmp/backend.log"
    exit 1
fi
echo "✅ Backend iniciado (PID: $BACKEND_PID)"

echo ""
echo "🚀 PASO 6: Instalando dependencias del frontend..."
cd /Users/paolacasadiegos/Downloads/truora/goland2026/frontEndGoland2026
npm install --legacy-peer-deps > /dev/null 2>&1
echo "✅ Dependencias instaladas"

echo ""
echo "🚀 PASO 7: Iniciando frontend en puerto 5173..."
npm run dev > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!
sleep 5
echo "✅ Frontend iniciado (PID: $FRONTEND_PID)"

echo ""
echo "================================"
echo "✅ SISTEMA INICIADO CORRECTAMENTE"
echo "================================"
echo ""
echo "📱 Frontend:  http://127.0.0.1:5173"
echo "🔧 Backend:   http://localhost:8081"
echo ""
echo "📝 Credenciales de prueba:"
echo "   Email:    test@example.com"
echo "   Password: password123"
echo ""
echo "📊 Datos de prueba: Mock data (10+ desafíos con paginación)"
echo ""
echo "⚠️  IMPORTANTE:"
echo "   - Para detener todo: bash STOP.sh"
echo "   - Backend logs: tail -f /tmp/backend.log"
echo "   - Frontend logs: tail -f /tmp/frontend.log"
echo ""
echo "================================"

