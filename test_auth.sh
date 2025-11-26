#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Testing NexusFlow Authentication System${NC}\n"

# Test 1: Register a new user
echo -e "${YELLOW}Test 1: Register New User${NC}"
REGISTER_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "test@nexusflow.io",
  "password": "password123",
  "display_name": "Test User"
}' localhost:50051 nexusflow.user.v1.AuthService/Register 2>&1)

if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    echo -e "${GREEN}✓ Registration successful!${NC}"
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token": "[^"]*"' | cut -d'"' -f4)
    echo "JWT Token: ${TOKEN:0:50}..."
else
    echo -e "${RED}✗ Registration failed${NC}"
    echo "$REGISTER_RESPONSE"
fi

echo ""

# Test 2: Login with the user
echo -e "${YELLOW}Test 2: Login${NC}"
LOGIN_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "test@nexusflow.io",
  "password": "password123"
}' localhost:50051 nexusflow.user.v1.AuthService/Login 2>&1)

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo -e "${GREEN}✓ Login successful!${NC}"
    echo "$LOGIN_RESPONSE" | grep -E '"(id|email|displayName)"'
else
    echo -e "${RED}✗ Login failed${NC}"
    echo "$LOGIN_RESPONSE"
fi

echo ""

# Test 3: Try to register with same email (should fail)
echo -e "${YELLOW}Test 3: Duplicate Email Registration (should fail)${NC}"
DUP_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "test@nexusflow.io",
  "password": "password456",
  "display_name": "Another User"
}' localhost:50051 nexusflow.user.v1.AuthService/Register 2>&1)

if echo "$DUP_RESPONSE" | grep -q "already exists"; then
    echo -e "${GREEN}✓ Correctly rejected duplicate email${NC}"
else
    echo -e "${RED}✗ Should have rejected duplicate email${NC}"
    echo "$DUP_RESPONSE"
fi

echo ""

# Test 4: Try login with wrong password (should fail)
echo -e "${YELLOW}Test 4: Wrong Password (should fail)${NC}"
WRONG_PW_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "test@nexusflow.io",
  "password": "wrongpassword"
}' localhost:50051 nexusflow.user.v1.AuthService/Login 2>&1)

if echo "$WRONG_PW_RESPONSE" | grep -q "invalid email or password"; then
    echo -e "${GREEN}✓ Correctly rejected wrong password${NC}"
else
    echo -e "${RED}✗ Should have rejected wrong password${NC}"
    echo "$WRONG_PW_RESPONSE"
fi

echo ""

# Test 5: Request password reset
echo -e "${YELLOW}Test 5: Request Password Reset${NC}"
RESET_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "test@nexusflow.io"
}' localhost:50051 nexusflow.user.v1.AuthService/RequestPasswordReset 2>&1)

if echo "$RESET_RESPONSE" | grep -q "success"; then
    echo -e "${GREEN}✓ Password reset requested${NC}"
    RESET_TOKEN=$(echo "$RESET_RESPONSE" | grep -o '"token": "[^"]*"' | cut -d'"' -f4)
    echo "Reset Token: ${RESET_TOKEN:0:50}..."
    
    # Test 6: Reset password with token
    echo ""
    echo -e "${YELLOW}Test 6: Reset Password${NC}"
    RESET_PW_RESPONSE=$(grpcurl -plaintext -d "{
      \"token\": \"$RESET_TOKEN\",
      \"new_password\": \"newpassword123\"
    }" localhost:50051 nexusflow.user.v1.AuthService/ResetPassword 2>&1)
    
    if echo "$RESET_PW_RESPONSE" | grep -q "success"; then
        echo -e "${GREEN}✓ Password reset successful${NC}"
        
        # Test 7: Login with new password
        echo ""
        echo -e "${YELLOW}Test 7: Login with New Password${NC}"
        NEW_LOGIN_RESPONSE=$(grpcurl -plaintext -d '{
          "email": "test@nexusflow.io",
          "password": "newpassword123"
        }' localhost:50051 nexusflow.user.v1.AuthService/Login 2>&1)
        
        if echo "$NEW_LOGIN_RESPONSE" | grep -q "token"; then
            echo -e "${GREEN}✓ Login with new password successful!${NC}"
        else
            echo -e "${RED}✗ Login with new password failed${NC}"
        fi
    else
        echo -e "${RED}✗ Password reset failed${NC}"
        echo "$RESET_PW_RESPONSE"
    fi
else
    echo -e "${RED}✗ Password reset request failed${NC}"
    echo "$RESET_RESPONSE"
fi

echo ""
echo -e "${YELLOW}✅ Testing Complete!${NC}"
